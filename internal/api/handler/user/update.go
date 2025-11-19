package user

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *User) PutApi1UsersId(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	data, err := request.ParseBody[dto.UpdateUserRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.user.UpdateByID(ctx, id, model.User{
		Name:  data.Name,
		Login: data.Login,
		Phone: data.Phone,
	})
	if validationErr := new(model.ValidationError); errors.As(err, &validationErr) {
		response.BadRequest(w, err.Error())
		return
	}
	if err != nil {
		response.BadRequest(w, "Пользователь с таким логином уже существует")
		return
	}

	response.OK(w, nil)
	return
}
