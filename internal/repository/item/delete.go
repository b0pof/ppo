package item

import (
	"context"
	"fmt"
)

func (r *Repository) DeleteByID(ctx context.Context, itemID int64) error {
	q := `delete from item where id = $1`

	_, err := r.db.ExecContext(ctx, q, itemID)
	if err != nil {
		return fmt.Errorf("item repository.DeleteByID: %w", err)
	}

	return nil
}
