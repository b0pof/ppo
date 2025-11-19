package order

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Order) PatchApi1UsersUserIdOrdersOrderId(w http.ResponseWriter, r *http.Request, _ int64, orderId int64) {
	ctx := r.Context()

	data, err := request.ParseBody[dto.UpdateOrderRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.order.UpdateStatus(ctx, orderId, data.Status)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Заказ не найден")
		return
	}
	if err != nil {
		h.log.Warn("failed to update order status", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, nil)
	return
}
