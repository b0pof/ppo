package review

import (
	"net/http"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	authUtil "github.com/b0pof/ppo/internal/util/auth"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Review) PostApi1ItemsIdReviews(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	userID := authUtil.GetUserID(ctx)

	data, err := request.ParseBody[dto.CreateReviewRequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	if id == 0 || data.Rating < 1 || data.Rating > 5 {
		response.BadRequest(w, "Item ID and valid rating (1-5) are required")
		return
	}

	newReview := model.Review{
		UserID:        userID,
		ItemID:        id,
		Rating:        data.Rating,
		Advantages:    data.Advantages,
		Disadvantages: data.Disadvantages,
		Comment:       data.Comment,
	}

	reviewID, err := h.review.AddReview(ctx, newReview)
	if err != nil {
		h.log.Warn("failed to create review", err)
		response.Internal(w, err)

		return
	}

	response.OK(w, dto.CreateReviewResponse{
		ReviewId: reviewID,
	})
}
