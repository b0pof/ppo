//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package item

import (
	"context"

	"github.com/b0pof/ppo/internal/model"
)

type itemRepo interface {
	Create(ctx context.Context, item model.Item) (int64, error)
	GetByID(ctx context.Context, itemID int64, userID int64) (model.ItemExtended, error)
	GetAllItems(ctx context.Context, userID int64) ([]model.ItemExtended, error)
	GetItemsBySellerID(ctx context.Context, sellerID int64) ([]model.Item, error)
	GetItemsByOrderID(ctx context.Context, orderID int64) ([]model.OrderItemInfo, error)
	DeleteByID(ctx context.Context, itemID int64) error
	UpdateByID(ctx context.Context, item model.Item) error
}
