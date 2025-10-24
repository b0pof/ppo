package user

import (
	"errors"
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *User) GetApi1UsersId(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	user, err := h.user.GetByID(ctx, id)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Пользователь не найден")
		return
	}
	if err != nil {
		h.log.Warn("failed to fetch user", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, toUserDTO(user))
	return
}

func toUserDTO(user model.User) dto.UserResponse {
	return dto.UserResponse{
		Name:  user.Name,
		Login: user.Login,
		Phone: user.Phone,
		Role:  user.Role,
	}
}
