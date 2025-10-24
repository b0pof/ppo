package order

import (
	"errors"
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Order) GetApi1UsersUserIdOrdersOrderId(w http.ResponseWriter, r *http.Request, _ int64, orderId int64) {
	ctx := r.Context()

	order, err := h.order.GetByID(ctx, orderId)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Заказ не найден")
		return
	}
	if err != nil {
		h.log.Warn("failed to fetch order", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, toGetOrderDTO(order))
	return
}

func toGetOrderDTO(o model.OrderExtended) dto.GetOrderResponse {
	items := make([]dto.OrderItemInfo, 0, len(o.Items))

	for _, i := range o.Items {
		items = append(items, dto.OrderItemInfo{
			Id:     i.ID,
			Price:  i.Price,
			Name:   i.ProductName,
			Count:  i.Count,
			ImgSrc: i.ImgSrc,
		})
	}

	return dto.GetOrderResponse{
		Id:         o.ID,
		Status:     o.Status,
		BuyerId:    o.BuyerID,
		ItemsCount: o.ItemsCount,
		CreatedAt:  o.CreatedAt,
		Sum:        o.Sum,
		Items:      items,
	}
}
