package cart

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) GetCartItemsAmount(ctx context.Context, userID int64) (int, error) {
	q := `select count(*) from cart_item where user_id = $1;`

	var count int

	err := r.db.GetContext(ctx, &count, q, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("repository.GetCartItemsAmount: %w", err)
	}

	return count, nil
}

func (r *Repository) GetCartContentByUserID(ctx context.Context, userID int64) (model.CartContent, error) {
	q := `
		select i.id, i.name, i.price, ci.count, i.imgsrc
		from item i
			join cart_item ci on i.id = ci.item_id and ci.user_id = $1;
	`

	var cartItems []cartItemRow

	err := r.db.SelectContext(ctx, &cartItems, q, userID)
	if err != nil {
		return model.CartContent{}, fmt.Errorf("item repository.GetCartItemsByUserID: %w", err)
	}

	//if len(cartItems) == 0 {
	//	return model.CartContent{}, model.ErrNotFound
	//}

	var totalPrice int
	var totalCount int

	for _, item := range cartItems {
		totalPrice += item.Price * item.Count
		totalCount += item.Count
	}

	return model.CartContent{
		TotalPrice: totalPrice,
		TotalCount: totalCount,
		Items:      convertCartItems(cartItems),
	}, nil
}
