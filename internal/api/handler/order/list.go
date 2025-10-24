package order

import (
	"errors"
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Order) GetApi1UsersIdOrders(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	orders, err := h.order.GetAllOrders(ctx, id)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Заказов не найдено")
		return
	}
	if err != nil {
		h.log.Warn("failed to fetch orders", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, toGetOrdersDTO(orders))
	return
}

func toGetOrdersDTO(orders []model.Order) dto.GetOrdersResponse {
	dtoOrders := make([]dto.OrderInfo, 0, len(orders))

	for _, o := range orders {
		dtoOrders = append(dtoOrders, dto.OrderInfo{
			BuyerId:    o.BuyerID,
			Id:         o.ID,
			ItemsCount: o.ItemsCount,
			Status:     o.Status,
			CreatedAt:  o.CreatedAt,
		})
	}

	return dto.GetOrdersResponse{
		Orders: dtoOrders,
	}
}
