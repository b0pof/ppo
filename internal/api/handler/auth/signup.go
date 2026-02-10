package auth

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"

	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Auth) PostApi1Users(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := request.ParseBody[dto.SignupRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.auth.Signup(ctx, data.Login, data.Password, string(data.Role))
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

	response.OK(w, "Код отправлен на почту")
}
