package auth

import (
	"errors"
	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"net/http"
	"time"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/request"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Auth) PostApi1Users(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := request.ParseBody[dto.SignupRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	sessionID, err := h.auth.Signup(ctx, data.Login, data.Password, string(data.Role))
	if errors.Is(err, model.ErrAlreadyExists) {
		response.BadRequest(w, "Пользователь уже существует")
		return
	}
	if errors.Is(err, model.ErrInvalidInput) {
		response.BadRequest(w, err.Error())
		return
	}
	if errors.Is(err, model.ErrWrongPassword) {
		response.BadRequest(w, err.Error())
		return
	}
	if err != nil {
		h.log.Warn("signup failed", err)
		response.Internal(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Expires:  time.Now().Add(h.sessionTTL),
		Path:     "/",
	}

	http.SetCookie(w, cookie)

	response.OK(w, sessionID)
	return
}
