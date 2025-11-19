package response_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/b0pof/ppo/internal/util/http/response"
)

type customResponseWriter struct {
	status int
	body   bytes.Buffer
}

func (c *customResponseWriter) Header() http.Header {
	return http.Header{}
}

func (c *customResponseWriter) WriteHeader(statusCode int) {
	c.status = statusCode
}

func (c *customResponseWriter) Write(b []byte) (int, error) {
	c.body.Write(b)

	return len(b), nil
}

type bodyData struct {
	Result string `json:"result"`
}

type errorResult struct {
	Error string `json:"error"`
}

func TestOK(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		w            *customResponseWriter
		data         interface{}
		expectations func(t assert.TestingT, w *customResponseWriter)
	}{
		{
			name: "success",
			w:    &customResponseWriter{},
			data: bodyData{"ok"},
			expectations: func(t assert.TestingT, w *customResponseWriter) {
				gotBody := w.body.String()
				expectedBody, _ := json.Marshal(bodyData{Result: "ok"})

				assert.Equal(t, string(expectedBody)+"\n", gotBody)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			OK(tc.w, tc.data)

			tc.expectations(t, tc.w)
		})
	}
}

func TestUnauthorized(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		w            *customResponseWriter
		data         interface{}
		expectations func(t assert.TestingT, w *customResponseWriter)
	}{
		{
			name: "success",
			w:    &customResponseWriter{},
			data: bodyData{"unathorized"},
			expectations: func(t assert.TestingT, w *customResponseWriter) {
				gotBody := w.body.String()
				expectedBody, _ := json.Marshal("Авторизация отсутствует")

				assert.Equal(t, string(expectedBody)+"\n", gotBody)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			Unauthorized(tc.w)

			tc.expectations(t, tc.w)
		})
	}
}

func TestBadRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		w            *customResponseWriter
		data         interface{}
		expectations func(t assert.TestingT, w *customResponseWriter)
	}{
		{
			name: "success",
			w:    &customResponseWriter{},
			data: bodyData{"bad request"},
			expectations: func(t assert.TestingT, w *customResponseWriter) {
				gotBody := w.body.String()
				expectedBody, _ := json.Marshal(errorResult{Error: "error"})

				assert.Equal(t, string(expectedBody)+"\n", gotBody)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			BadRequest(tc.w, "error")

			tc.expectations(t, tc.w)
		})
	}
}

func TestForbidden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		w            *customResponseWriter
		data         interface{}
		expectations func(t assert.TestingT, w *customResponseWriter)
	}{
		{
			name: "success",
			w:    &customResponseWriter{},
			data: bodyData{"forbidden"},
			expectations: func(t assert.TestingT, w *customResponseWriter) {
				gotBody := w.body.String()

				assert.Equal(t, "", gotBody)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			Forbidden(tc.w, nil)

			tc.expectations(t, tc.w)
		})
	}
}
