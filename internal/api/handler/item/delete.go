package item

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Item) DeleteApi1ItemsId(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	err := h.item.Delete(ctx, id)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Товар не найден")
		return
	}
	if err != nil {
		h.log.Warn("failed to delete item", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, nil)
	return
}
