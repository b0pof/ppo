package order

import (
	"errors"
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/request"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
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
