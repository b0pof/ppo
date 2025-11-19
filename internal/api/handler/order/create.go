package order

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Order) PostApi1UsersIdOrders(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	orderID, err := h.order.Create(ctx, id)
	if errors.Is(err, model.ErrCartIsEmpty) {
		response.BadRequest(w, "Корзина пуста")
		return
	}
	if err != nil {
		h.log.Warn("failed to create order", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, dto.CreateOrderResponse{OrderID: orderID})
	return
}
