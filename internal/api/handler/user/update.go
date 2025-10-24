package user

import (
	"errors"
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/request"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
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
