package item

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	authUtil "github.com/b0pof/ppo/internal/util/auth"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Item) PostApi1Items(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sellerID := authUtil.GetUserID(ctx)

	data, err := request.ParseBody[dto.CreateItemRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	itemID, err := h.item.Create(ctx, model.Item{
		Name:        data.Name,
		Description: data.Description,
		Seller: model.Seller{
			ID: sellerID,
		},
		Price:  data.Price,
		ImgSrc: data.ImgSrc,
	})
	if validationErr := new(model.ValidationError); errors.As(err, &validationErr) {
		response.BadRequest(w, err.Error())
		return
	}
	if err != nil {
		h.log.Warn("failed to create item", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, dto.CreateItemResponse{
		ItemId: itemID,
	})
	return
}
