package item

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	authUtil "github.com/b0pof/ppo/internal/util/auth"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Item) GetApi1ItemsId(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	item, err := h.item.GetByID(ctx, authUtil.GetUserID(ctx), id)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Товар не найден")
		return
	}
	if err != nil {
		h.log.Warn("failed to fetch all items", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, itemToDTO(item))
	return
}

func itemToDTO(item model.ItemExtended) dto.Item {
	return dto.Item{
		Id:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		ImgSrc:      item.ImgSrc,
		Rating:      item.Rating,
		Seller:      item.Seller.Name,
		Amount:      item.Amount,
	}
}
