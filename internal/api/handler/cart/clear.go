package cart

import (
	"net/http"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Cart) DeleteApi1UsersIdCartItems(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	err := h.cart.Clear(ctx, id)
	if err != nil {
		h.log.Warn("failed to clear cart", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, nil)
	return
}
