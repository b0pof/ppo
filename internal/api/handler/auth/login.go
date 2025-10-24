package auth

import (
	"errors"
	"net/http"
	"time"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/request"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Auth) PostApi1Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := request.ParseBody[dto.LoginRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	sessionID, err := h.auth.Login(ctx, data.Login, data.Password)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Пользователь не найден")
		return
	}
	if errors.Is(err, model.ErrInvalidInput) {
		response.BadRequest(w, err.Error())
		return
	}
	if errors.Is(err, model.ErrWrongPassword) {
		response.BadRequest(w, "Неверный пароль")
		return
	}
	if err != nil {
		h.log.Warn("login failed", err)
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
