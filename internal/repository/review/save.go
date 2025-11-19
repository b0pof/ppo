package review

import (
	"context"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) AddReview(ctx context.Context, review model.Review) (int64, error) {
	q := `
        INSERT INTO review (
            user_id,
            item_id,
            rating,
            advantages,
            disadvantages,
            comment
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	var reviewID int64

	err := r.db.GetContext(ctx, &reviewID, q, review.UserID, review.ItemID, review.Rating, review.Advantages, review.Disadvantages, review.Comment)
	if err != nil {
		return 0, fmt.Errorf("repository.AddReview: insert failed: %w", err)
	}

	return reviewID, nil
}
