package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) GetReviews(ctx context.Context, itemID int64) ([]model.Review, error) {
	q := `
        SELECT 
            r.id,
            r.user_id,
            r.rating,
            r.advantages,
            r.disadvantages,
            r.comment,
            r.created_at,
            u.name as user_name
        FROM review r
        JOIN "user" u ON r.user_id = u.id
        WHERE r.item_id = $1
        ORDER BY r.created_at DESC
    `

	var reviews []struct {
		ID            int64     `db:"id"`
		UserID        int64     `db:"user_id"`
		Rating        int       `db:"rating"`
		Advantages    string    `db:"advantages"`
		Disadvantages string    `db:"disadvantages"`
		Comment       string    `db:"comment"`
		CreatedAt     time.Time `db:"created_at"`
		Username      string    `db:"user_name"`
	}

	err := r.db.SelectContext(ctx, &reviews, q, itemID)
	if errors.Is(err, sql.ErrNoRows) {
		return []model.Review{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.GetReviews: %w", err)
	}

	result := make([]model.Review, 0, len(reviews))
	for _, review := range reviews {
		result = append(result, model.Review{
			ID:            review.ID,
			UserID:        review.UserID,
			Username:      review.Username,
			Rating:        review.Rating,
			Advantages:    review.Advantages,
			Disadvantages: review.Disadvantages,
			Comment:       review.Comment,
			CreatedAt:     review.CreatedAt,
		})
	}

	return result, nil
}
