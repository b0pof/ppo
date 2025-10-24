package cart

import (
	"errors"
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Cart) DeleteApi1UsersUserIdCartItemsItemId(w http.ResponseWriter, r *http.Request, userId int64, itemId int64) {
	ctx := r.Context()

	newCount, err := h.cart.DeleteItem(ctx, userId, itemId)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Товар не найден в корзине")
		return
	}
	if err != nil {
		h.log.Warn("failed to delete item from cart", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, dto.DeleteItemResponse{Count: newCount})
	return
}
