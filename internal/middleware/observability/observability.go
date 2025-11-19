package observability

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"

	"github.com/b0pof/ppo/internal/pkg/metrics"
)

var ErrHijackAssertion = errors.New("type assertion to http.Hijacker failed")

var currRequestID uint64

type responseWriterInterceptor struct {
	w     http.ResponseWriter
	wCopy bytes.Buffer
}

func newResponseWriterInterceptor(w http.ResponseWriter) *responseWriterInterceptor {
	var wCopy bytes.Buffer
	return &responseWriterInterceptor{
		w:     w,
		wCopy: wCopy,
	}
}

func (wi *responseWriterInterceptor) WriteHeader(statusCode int) {
	wi.w.WriteHeader(statusCode)
}

func (wi *responseWriterInterceptor) Header() http.Header {
	return wi.w.Header()
}

func (wi *responseWriterInterceptor) Write(d []byte) (int, error) {
	wi.wCopy.Write(d)
	return wi.w.Write(d)
}

func (wi *responseWriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := wi.w.(http.Hijacker)
	if !ok {
		return nil, nil, ErrHijackAssertion
	}
	return h.Hijack()
}

func (wi *responseWriterInterceptor) GetStatusCode() (int, error) {
	type response struct {
		StatusCode int `json:"status"`
	}
	var resp response
	if err := json.Unmarshal(wi.wCopy.Bytes(), &resp); err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func New(metrics metrics.Collector, l *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddUint64(&currRequestID, 1)
			requestLogger := l.With(slog.Uint64("requestID", currRequestID))
			requestLogger.Info("new",
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI))

			wi := newResponseWriterInterceptor(w)

			start := time.Now()
			next.ServeHTTP(wi, r)
			dur := time.Since(start)

			statusCode, err := wi.GetStatusCode()
			if err != nil {
				requestLogger.Error("error while getting status code",
					slog.String("duration", dur.String()))
				return
			}
			requestLogger.Info("response",
				slog.Int("statusCode", statusCode),
				slog.String("duration", dur.String()))

			if statusCode >= 300 {
				metrics.IncreaseErr(strconv.Itoa(statusCode), r.RequestURI)
			}
			path := r.URL.Path
			pathVars := mux.Vars(r)
			for key, value := range pathVars {
				path, _ = strings.CutSuffix(path, value)
				path += fmt.Sprintf("{%s}", key)
			}
			metrics.AddDurationToHistogram(path, dur)
			metrics.AddDurationToSummary(strconv.Itoa(statusCode), path, dur)
			metrics.IncreaseHits(path)
		})
	}
}
