//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package order

import (
	"context"

	"github.com/b0pof/ppo/internal/model"
)

type orderRepo interface {
	Create(ctx context.Context, userID int64, items []model.CartItem) (int64, error)
	GetByID(ctx context.Context, orderID int64) (model.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, newStatus string) error
}

type itemRepo interface {
	GetItemsByOrderID(ctx context.Context, orderID int64) ([]model.OrderItemInfo, error)
}

type cartRepo interface {
	GetCartContentByUserID(ctx context.Context, userID int64) (model.CartContent, error)
	DeleteAllCartItems(ctx context.Context, userID int64) error
}
