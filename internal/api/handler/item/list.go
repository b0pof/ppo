package item

import (
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	authUtil "github.com/b0pof/ppo/internal/util/auth"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Item) GetApi1Items(w http.ResponseWriter, r *http.Request, _ dto.GetApi1ItemsParams) {
	ctx := r.Context()

	items, err := h.item.GetAllItems(ctx, authUtil.GetUserID(ctx))
	if err != nil {
		h.log.Warn("failed to fetch all items", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, dto.ItemsFetchResponse{ProductCards: itemsToDTO(items)})
	return
}

func itemsToDTO(cards []model.ItemExtended) []dto.ItemCard {
	var productCards []dto.ItemCard

	for _, card := range cards {
		productCards = append(productCards, dto.ItemCard{
			Id:     card.ID,
			Name:   card.Name,
			Price:  card.Price,
			ImgSrc: card.ImgSrc,
			Seller: card.Seller.Name,
			Amount: &card.Amount,
		})
	}

	return productCards
}
