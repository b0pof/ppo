package cookie

import (
	"net/http"
	"time"
)

const (
	sessionKey = "session_id"
	sessionTTL = 12 * time.Hour
)

func GetSession(r *http.Request) (string, error) {
	sessionID, err := r.Cookie(sessionKey)
	if err != nil {
		return "", err
	}

	return sessionID.Value, nil
}

func SetSession(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     sessionKey,
		Value:    sessionID,
		HttpOnly: true,
		Expires:  time.Now().Add(sessionTTL),
		Path:     "/",
	}

	http.SetCookie(w, cookie)
}
