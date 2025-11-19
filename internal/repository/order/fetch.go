package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) GetByID(ctx context.Context, orderID int64) (model.Order, error) {
	q := `
		select o.id, o.buyer_id, sum(oi.count) as count, sum(i.price * oi.count) as sum, o.status, o.created_at
		from "order" o
			inner join order_item oi on o.id = oi.order_id
			inner join item i on oi.item_id = i.id
		where o.id = $1
		group by o.id, o.buyer_id, o.status, o.created_at;
	`

	var order orderRow

	err := r.db.GetContext(ctx, &order, q, orderID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Order{}, model.ErrNotFound
	}
	if err != nil {
		return model.Order{}, fmt.Errorf("order repository.GetByID: %w", err)
	}

	return convertOrder(order), nil
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error) {
	q := `
		select o.id, o.buyer_id, sum(oi.count) as count, sum(i.price * oi.count) as sum, o.status, o.created_at
		from "order" o
			inner join order_item oi on o.id = oi.order_id
			inner join item i on oi.item_id = i.id
		where o.buyer_id = $1
		group by o.id, o.buyer_id, o.status, o.created_at;
	`

	var orders []orderRow

	err := r.db.SelectContext(ctx, &orders, q, userID)
	if err != nil {
		return nil, fmt.Errorf("order repository.GetOrdersByUserID: %w", err)
	}

	if len(orders) == 0 {
		return nil, model.ErrNotFound
	}

	return convertOrders(orders), nil
}
