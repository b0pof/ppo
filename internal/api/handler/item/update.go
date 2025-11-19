package item

import (
	"errors"
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	authUtil "github.com/b0pof/ppo/internal/util/auth"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
	"github.com/b0pof/ppo/internal/util/pointer"
)

func (h *Item) PutApi1ItemsId(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	sellerID := authUtil.GetUserID(ctx)

	data, err := request.ParseBody[dto.UpdateItemRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.item.Update(ctx, model.Item{
		ID:          id,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Seller: model.Seller{
			ID: sellerID,
		},
		ImgSrc: data.ImgSrc,
	})
	if validationErr := new(model.ValidationError); errors.As(err, &validationErr) {
		response.BadRequest(w, err.Error())
		return
	}
	if errors.Is(err, model.ErrNoAccess) {
		response.Forbidden(w, pointer.To("Доступ к редактированию товара ограничен"))
		return
	}
	if err != nil {
		h.log.Warn("failed to create item", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, nil)
	return
}
