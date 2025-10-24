package item

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
)

const (
	pgStringField = "'%s'"
)

func (r *Repository) UpdateByID(ctx context.Context, item model.Item) error {
	q := `
		update item set
		    name = %s,
			price = %s,
		    description = %s,
			imgsrc = %s
		where id = $1;
	`

	newName := "name"
	if item.Name != "" {
		newName = fmt.Sprintf(pgStringField, item.Name)
	}

	newPrice := "price"
	if item.Price != 0 {
		newPrice = fmt.Sprintf("%d", item.Price)
	}

	newDescription := "description"
	if item.Description != "" {
		newDescription = fmt.Sprintf(pgStringField, item.Description)
	}

	newImgSrc := "imgsrc"
	if item.ImgSrc != "" {
		newImgSrc = fmt.Sprintf(pgStringField, item.ImgSrc)
	}

	q = fmt.Sprintf(q, newName, newPrice, newDescription, newImgSrc)

	_, err := r.db.ExecContext(ctx, q, item.ID)
	if err != nil {
		return fmt.Errorf("item repository.UpdateByID: %w", err)
	}

	return nil
}
