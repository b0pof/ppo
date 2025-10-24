package review

import (
	"net/http"

	dto "git.iu7.bmstu.ru/kia22u475/ppo/internal/generated"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func (h *Review) GetApi1ItemsIdReviews(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	reviews, err := h.review.GetReviews(ctx, id)
	if err != nil {
		h.log.Warn("failed to get reviews", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, toReviewsDTO(reviews))
}

func toReviewsDTO(reviews []model.Review) dto.GetReviewsResponse {
	dtoReviews := make([]dto.ReviewInfo, 0, len(reviews))

	for _, r := range reviews {
		dtoReviews = append(dtoReviews, dto.ReviewInfo{
			Id:            r.ID,
			UserId:        r.UserID,
			UserName:      r.Username,
			Rating:        r.Rating,
			Advantages:    r.Advantages,
			Disadvantages: r.Disadvantages,
			Comment:       r.Comment,
			CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return dto.GetReviewsResponse{
		Reviews: dtoReviews,
	}
}
