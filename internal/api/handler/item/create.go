package item

import (
	"errors"
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	authUtil "git.iu7.bmstu.ru/kia22u475/ppo/internal/util/auth"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/request"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
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
