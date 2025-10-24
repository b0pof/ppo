package order

import (
	"context"
	"fmt"
)

func (r *Repository) UpdateStatus(ctx context.Context, orderID int64, newStatus string) error {
	q := `update "order" set status = $1 where id = $2`

	_, err := r.db.ExecContext(ctx, q, newStatus, orderID)
	if err != nil {
		return fmt.Errorf("order repository.UpdateStatus: %w", err)
	}

	return nil
}
