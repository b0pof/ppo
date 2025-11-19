package cart

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) DeleteCartItem(ctx context.Context, userID int64, itemID int64) (int, error) {
	q := `
		update cart_item set count = cart_item.count - 1
		where user_id = $1 AND item_id = $2
		returning cart_item.count AS count;
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("repository.DeleteCartItem: failed to begin transaction: %w", err)
	}

	var count int

	err = tx.GetContext(ctx, &count, q, userID, itemID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, model.ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("repository.DeleteCartItem: failed to update: %w", err)
	}

	if count == 0 {
		q = `delete from cart_item where user_id = $1 AND item_id = $2;`

		if _, err = tx.ExecContext(ctx, q, userID, itemID); err != nil {
			return 0, fmt.Errorf("repository.DeleteCartItem: failed to delete: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("repository.DeleteCartItem: failed to commit transaction: %w", err)
	}

	return count, nil
}

func (r *Repository) DeleteAllCartItems(ctx context.Context, userID int64) error {
	q := `delete from cart_item where user_id = $1;`

	_, err := r.db.ExecContext(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("repository.DeleteAllCartItems: %w", err)
	}

	return nil
}
