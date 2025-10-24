package order

import (
	"time"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
)

type orderRow struct {
	ID        int64     `db:"id"`
	BuyerID   int64     `db:"buyer_id"`
	Count     int       `db:"count"`
	Sum       int       `db:"sum"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

func convertOrder(row orderRow) model.Order {
	return model.Order{
		ID:         row.ID,
		BuyerID:    row.BuyerID,
		Sum:        row.Sum,
		Status:     row.Status,
		CreatedAt:  row.CreatedAt,
		ItemsCount: row.Count,
	}
}

func convertOrders(rows []orderRow) []model.Order {
	orders := make([]model.Order, 0, len(rows))

	for _, row := range rows {
		orders = append(orders, convertOrder(row))
	}

	return orders
}
