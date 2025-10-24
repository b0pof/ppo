package cart

import (
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Cart) GetApi1UsersIdCartItems(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	cart, err := h.cart.GetCartContent(ctx, id)
	if err != nil {
		h.log.Warn("failed to fetch cart content", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, cartContentToDTO(cart))
	return
}

func cartContentToDTO(cart model.CartContent) dto.GetCartItemsResponse {
	items := make([]dto.CartItem, 0, len(cart.Items))

	for _, i := range cart.Items {
		items = append(items, cartItemToDTO(i))
	}

	return dto.GetCartItemsResponse{
		TotalPrice: cart.TotalPrice,
		TotalCount: cart.TotalCount,
		Items:      items,
	}
}

func cartItemToDTO(item model.CartItem) dto.CartItem {
	return dto.CartItem{
		Id:     item.ID,
		Name:   item.Name,
		Price:  item.Price,
		Count:  item.Count,
		ImgSrc: item.ImgSrc,
	}
}
