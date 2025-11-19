package seller

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Seller) GetApi1SellerIdItems(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	items, err := h.item.GetItemsBySellerID(ctx, id)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Товары не найдены")
		return
	}
	if err != nil {
		h.log.Warn("failed to fetch seller items", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, dto.ItemsFetchResponse{ProductCards: itemsToDTO(items)})
	return
}

func itemsToDTO(items []model.Item) []dto.ItemCard {
	var itemCards []dto.ItemCard

	for _, item := range items {
		itemCards = append(itemCards, dto.ItemCard{
			Id:     item.ID,
			Name:   item.Name,
			Price:  item.Price,
			Seller: item.Seller.Name,
			ImgSrc: item.ImgSrc,
		})
	}

	return itemCards
}
