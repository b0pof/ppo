package item

import (
	"context"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) Create(ctx context.Context, item model.Item) (int64, error) {
	q := `
		insert into item (name, seller_id, price, description, imgsrc)
		values ($1, $2, $3, $4, $5)
		returning id;
	`

	var id int64

	err := r.db.GetContext(ctx, &id, q,
		item.Name,
		item.Seller.ID,
		item.Price,
		item.Description,
		item.ImgSrc,
	)

	if err != nil {
		return 0, fmt.Errorf("item repository.Create: %w", err)
	}

	return id, nil
}
