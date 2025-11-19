//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package cart

import (
	"context"

	"github.com/b0pof/ppo/internal/model"
)

type cartRepo interface {
	GetCartItemsAmount(ctx context.Context, userID int64) (int, error)
	GetCartContentByUserID(ctx context.Context, userID int64) (model.CartContent, error)
	AddItemToCart(ctx context.Context, userID int64, itemID int64) (int, error)
	DeleteCartItem(ctx context.Context, userID int64, itemID int64) (int, error)
	DeleteAllCartItems(ctx context.Context, userID int64) error
}
