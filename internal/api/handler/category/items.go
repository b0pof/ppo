package category

import (
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	authUtil "github.com/b0pof/ppo/internal/util/auth"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Category) GetApi1CategoriesIdItems(w http.ResponseWriter, r *http.Request, id int64, _ dto.GetApi1CategoriesIdItemsParams) {
	ctx := r.Context()
	userID := authUtil.GetUserID(ctx)

	items, err := h.item.GetItemsByCategoryID(ctx, id, userID)
	if err != nil {
		h.log.Warn("failed to get items by category", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, toItemsDTO(items))
}

func toItemsDTO(items []model.ItemExtended) dto.GetItemsByCategoryResponse {
	dtoItems := make([]dto.ItemInfo, 0, len(items))

	for _, i := range items {
		dtoItems = append(dtoItems, dto.ItemInfo{
			Id:          i.ID,
			Name:        i.Name,
			Description: i.Description,
			Price:       i.Price,
			ImgSrc:      i.ImgSrc,
			Seller: struct {
				Id   int64  `json:"id"`
				Name string `json:"name"`
			}{
				Id:   i.Seller.ID,
				Name: i.SellerName,
			},
			Amount: i.Amount,
		})
	}

	return dto.GetItemsByCategoryResponse{
		Items: dtoItems,
	}
}
