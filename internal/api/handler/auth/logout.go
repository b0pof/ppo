package auth

import (
	"errors"
	"net/http"
	"time"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Auth) DeleteApi1Auth(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) || session.Value == "" {
		response.Unauthorized(w)
		return
	}

	err = h.auth.Logout(session.Value)
	if errors.Is(err, model.ErrNotFound) {
		response.Forbidden(w, nil)
		return
	}
	if err != nil {
		h.log.Warn("logout failed", err)
		response.Internal(w, err)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	response.OK(w, nil)
}
