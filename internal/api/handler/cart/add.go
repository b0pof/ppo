package cart

import (
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Cart) PostApi1UsersIdCartItems(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	data, err := request.ParseBody[dto.AddItemRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	newCount, err := h.cart.AddItem(ctx, id, data.ItemId)
	if err != nil {
		h.log.Warn("failed to add item to cart", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, dto.AddItemResponse{Count: newCount})
	return
}
