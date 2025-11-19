package item

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) GetByID(ctx context.Context, itemID int64, userID int64) (model.ItemExtended, error) {
	q := `
		select i.id, i.name, i.seller_id, i.rating, s.name as seller_name, i.price, i.description, i.imgsrc,
		       coalesce(ci.count, 0) AS amount
		from item i
			join "user" s on i.seller_id = s.id
			left join cart_item ci on i.id = ci.item_id and ci.user_id = $1
		where i.id = $2;
	`

	var item itemExtendedRow

	err := r.db.GetContext(ctx, &item, q, userID, itemID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ItemExtended{}, model.ErrNotFound
	}
	if err != nil {
		return model.ItemExtended{}, fmt.Errorf("item repository.GetByID: %w", err)
	}

	return convertItemExtended(item), nil
}

func (r *Repository) GetAllItems(ctx context.Context, userID int64) ([]model.ItemExtended, error) {
	q := `
		select i.id, i.name, i.seller_id, s.name as seller_name, i.price, i.description, i.imgsrc,
			   coalesce(ci.count, 0) AS amount

		from item i
			join "user" s on i.seller_id = s.id
			left join cart_item ci on i.id = ci.item_id and ci.user_id = $1;
	`

	var items []itemExtendedRow

	err := r.db.SelectContext(ctx, &items, q, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("item repository.GetAllItems: %w", err)
	}

	return convertItemsExtended(items), nil
}

func (r *Repository) GetItemsBySellerID(ctx context.Context, sellerID int64) ([]model.Item, error) {
	q := `
		select i.id, i.name, i.description, i.price, i.imgsrc, i.seller_id, s.name as seller_name
		from item i
			join "user" s on i.seller_id = s.id
		where i.seller_id = $1;
	`

	var items []itemRow

	err := r.db.SelectContext(ctx, &items, q, sellerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("item repository.GetItemsBySellerID: %w", err)
	}

	return convertItems(items), nil
}

func (r *Repository) GetItemsByOrderID(ctx context.Context, orderID int64) ([]model.OrderItemInfo, error) {
	q := `
		select i.id, i.name, i.price, i.imgsrc, oi.count
		from item i
			join order_item oi on oi.item_id = i.id and oi.order_id = $1;
	`

	var items []orderItemRow

	err := r.db.SelectContext(ctx, &items, q, orderID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("item repository.GetItemsByOrderID: %w", err)
	}

	return convertOrderItems(items), nil
}

func (r *Repository) GetItemsByCategoryID(ctx context.Context, categoryID int64, userID int64) ([]model.ItemExtended, error) {
	q := `
        SELECT 
            i.id, 
            i.name, 
            i.seller_id, 
            s.name as seller_name, 
            i.price, 
            i.description, 
            i.imgsrc,
            COALESCE(ci.count, 0) AS amount
        FROM item i
        JOIN "user" s ON i.seller_id = s.id
        LEFT JOIN cart_item ci ON i.id = ci.item_id AND ci.user_id = $1
        WHERE EXISTS (
            SELECT 1 
            FROM item_category ic 
            WHERE ic.item_id = i.id AND ic.category_id = $2
        );
    `

	var items []itemExtendedRow

	err := r.db.SelectContext(ctx, &items, q, userID, categoryID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("item repository.GetItemsByCategoryID: %w", err)
	}

	return convertItemsExtended(items), nil
}
