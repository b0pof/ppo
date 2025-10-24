package auth

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Auth) GetApi1Auth(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) || session.Value == "" {
		response.Unauthorized(w)
		return
	}
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	if !h.auth.IsLoggedIn(session.Value) {
		response.Unauthorized(w)
		return
	}

	response.OK(w, nil)
	return
}
