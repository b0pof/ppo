package order

import (
	"context"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

type IOrderUsecase interface {
	Create(ctx context.Context, userID int64) (int64, error)
	GetByID(ctx context.Context, orderID int64) (model.OrderExtended, error)
	GetAllOrders(ctx context.Context, userID int64) ([]model.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, newStatus string) error
	Cancel(ctx context.Context, orderID int64) error
}

type Usecase struct {
	order orderRepo
	item  itemRepo
	cart  cartRepo
}

func New(o orderRepo, i itemRepo, c cartRepo) *Usecase {
	return &Usecase{
		order: o,
		item:  i,
		cart:  c,
	}
}

func (u *Usecase) Create(ctx context.Context, userID int64) (int64, error) {
	cart, err := u.cart.GetCartContentByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("order usecase: failed to get cart items: %w", err)
	}

	if len(cart.Items) == 0 {
		return 0, model.ErrCartIsEmpty
	}

	orderID, err := u.order.Create(ctx, userID, cart.Items)
	if err != nil {
		return 0, fmt.Errorf("order usecase: failed to create order: %w", err)
	}

	err = u.cart.DeleteAllCartItems(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("order usecase: failed to clear cart: %w", err)
	}

	return orderID, nil
}

func (u *Usecase) GetByID(ctx context.Context, orderID int64) (model.OrderExtended, error) {
	order, err := u.order.GetByID(ctx, orderID)
	if err != nil {
		return model.OrderExtended{}, fmt.Errorf("order usecase: failed to get order by ID: %w", err)
	}

	items, err := u.item.GetItemsByOrderID(ctx, order.ID)
	if err != nil {
		return model.OrderExtended{}, fmt.Errorf("order usecase: failed to get order items: %w", err)
	}

	var sum int
	for _, item := range items {
		sum += item.Price * item.Count
	}

	orderExtended := model.OrderExtended{
		Order: order,
		Sum:   sum,
		Items: items,
	}

	return orderExtended, nil
}

func (u *Usecase) GetAllOrders(ctx context.Context, userID int64) ([]model.Order, error) {
	orders, err := u.order.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("order usecase: failed to get all user orders: %w", err)
	}

	return orders, nil
}

func (u *Usecase) UpdateStatus(ctx context.Context, orderID int64, newStatus string) error {
	order, err := u.order.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order usecase: failed to get order by ID: %w", err)
	}

	err = u.order.UpdateStatus(ctx, order.ID, newStatus)
	if err != nil {
		return fmt.Errorf("order usecase: failed to update status: %w", err)
	}

	return nil
}

func (u *Usecase) Cancel(ctx context.Context, orderID int64) error {
	order, err := u.order.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order usecase: failed to get order by ID: %w", err)
	}

	err = u.order.UpdateStatus(ctx, order.ID, model.OrderStatusCancelled)
	if err != nil {
		return fmt.Errorf("order usecase: failed to update status: %w", err)
	}

	return nil
}
