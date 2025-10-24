package cart

import (
	"context"
	"fmt"
)

func (r *Repository) AddItemToCart(ctx context.Context, userID int64, itemID int64) (int, error) {
	q := `
		insert into cart_item(user_id, item_id) values ($1, $2)
		on conflict (user_id, item_id)
		do update set count = cart_item.count + 1
		returning cart_item.count AS count;
	`

	var count int

	err := r.db.GetContext(ctx, &count, q, userID, itemID)
	if err != nil {
		return 0, fmt.Errorf("repository.AddItemToCart: %w", err)
	}

	return count, nil
}
