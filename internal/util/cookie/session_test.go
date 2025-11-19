package cookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/b0pof/ppo/internal/util/cookie"
)

func TestGetSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		r            http.Request
		expectations func(t assert.TestingT, got string, err error)
	}{
		{
			name: "success",
			r: http.Request{
				Header: http.Header{
					"Cookie": []string{"session_id=session_hash;"},
				},
			},
			expectations: func(t assert.TestingT, got string, err error) {
				assert.Equal(t, "session_hash", got)
				assert.NoError(t, err)
			},
		},
		{
			name: "no cookie",
			r:    http.Request{},
			expectations: func(t assert.TestingT, got string, err error) {
				assert.ErrorIs(t, err, http.ErrNoCookie)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := GetSession(&tc.r)

			tc.expectations(t, out, err)
		})
	}
}

func TestSetSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		w            http.ResponseWriter
		sessionID    string
		expectations func(t assert.TestingT, w http.ResponseWriter)
	}{
		{
			name:      "success",
			w:         httptest.NewRecorder(),
			sessionID: "session_hash",
			expectations: func(t assert.TestingT, w http.ResponseWriter) {
				cookie := w.Header().Get("Set-Cookie")
				parsedCookie, _ := http.ParseSetCookie(cookie)

				assert.Equal(t, parsedCookie.Value, "session_hash")
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			SetSession(tc.w, tc.sessionID)

			tc.expectations(t, tc.w)
		})
	}
}
