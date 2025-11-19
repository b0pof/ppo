package order

import (
	"context"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) Create(ctx context.Context, userID int64, items []model.CartItem) (int64, error) {
	q := `
		insert into "order" (buyer_id, status)
		values ($1, $2)
		returning id;
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("order repository.Create: failed to begin transaction: %w", err)
	}

	var orderID int64

	err = tx.GetContext(ctx, &orderID, q, userID, model.OrderStatusCreated)
	if err != nil {
		return 0, fmt.Errorf("order repository.Create: failed to create order: %w", err)
	}

	q = `
		insert into order_item (order_id, item_id, count)
		values (:order_id, :item_id, :count);
	`

	args := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		args = append(args, map[string]interface{}{
			"order_id": orderID,
			"item_id":  item.ID,
			"count":    item.Count,
		})
	}

	_, err = tx.NamedExecContext(ctx, q, args)
	if err != nil {
		return 0, fmt.Errorf("order repository.Create: failed to insert items: %w", err)
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("order repository.Create: failed to commit transaction: %w", err)
	}

	return orderID, nil
}
