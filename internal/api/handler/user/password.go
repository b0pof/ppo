package user

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *User) PatchApi1UsersIdPassword(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	data, err := request.ParseBody[dto.UpdatePasswordRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.user.UpdatePassword(ctx, id, data.Password, data.NewPassword)
	if errors.Is(err, model.ErrWrongPassword) {
		response.BadRequest(w, "Неверный пароль")
		return
	}
	if validationErr := new(model.ValidationError); errors.As(err, &validationErr) {
		response.BadRequest(w, err.Error())
		return
	}
	if err != nil {
		h.log.Warn("failed to update password", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, nil)
	return
}
