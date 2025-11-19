package user

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *User) GetApi1UsersIdMeta(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	user, err := h.user.GetMetaInfoByUserID(ctx, id)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Пользователь не найден")
		return
	}
	if err != nil {
		h.log.Warn("failed to fetch user meta", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, toUserMetaDTO(user))
	return
}

func toUserMetaDTO(user model.UserMetaInfo) dto.UserMetaResponse {
	return dto.UserMetaResponse{
		Name:            user.Name,
		CartItemsAmount: user.CartItemsAmount,
	}
}
