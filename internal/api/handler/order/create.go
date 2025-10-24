package order

import (
	"errors"
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
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
