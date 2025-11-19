package cart

import (
	"context"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

type ICartUsecase interface {
	GetCartItemsAmount(ctx context.Context, userID int64) (int, error)
	GetCartContent(ctx context.Context, userID int64) (model.CartContent, error)
	AddItem(ctx context.Context, userID int64, itemID int64) (int, error)
	DeleteItem(ctx context.Context, userID int64, itemID int64) (int, error)
	Clear(ctx context.Context, userID int64) error
}

type Usecase struct {
	cart cartRepo
}

func New(r cartRepo) *Usecase {
	return &Usecase{
		cart: r,
	}
}

func (u *Usecase) GetCartItemsAmount(ctx context.Context, userID int64) (int, error) {
	amount, err := u.cart.GetCartItemsAmount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("cart usecase: failed to get cart items amount: %w", err)
	}

	return amount, nil
}

func (u *Usecase) GetCartContent(ctx context.Context, userID int64) (model.CartContent, error) {
	items, err := u.cart.GetCartContentByUserID(ctx, userID)
	if err != nil {
		return model.CartContent{}, fmt.Errorf("item usecase: failed to get cart items: %w", err)
	}

	return items, nil
}

func (u *Usecase) AddItem(ctx context.Context, userID int64, itemID int64) (int, error) {
	newAmount, err := u.cart.AddItemToCart(ctx, userID, itemID)
	if err != nil {
		return 0, fmt.Errorf("cart usecase: failed to add item to cart: %w", err)
	}

	return newAmount, nil
}

func (u *Usecase) DeleteItem(ctx context.Context, userID int64, itemID int64) (int, error) {
	newAmount, err := u.cart.DeleteCartItem(ctx, userID, itemID)
	if err != nil {
		return 0, fmt.Errorf("cart usecase: failed to delete cart item: %w", err)
	}

	return newAmount, nil
}

func (u *Usecase) Clear(ctx context.Context, userID int64) error {
	err := u.cart.DeleteAllCartItems(ctx, userID)
	if err != nil {
		return fmt.Errorf("cart usecase: failed to delete cart items: %w", err)
	}

	return nil
}
