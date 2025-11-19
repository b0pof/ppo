package item

import (
	"context"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

type IItemUsecase interface {
	Create(ctx context.Context, item model.Item) (int64, error)
	GetByID(ctx context.Context, userID int64, itemID int64) (model.ItemExtended, error)
	GetAllItems(ctx context.Context, userID int64) ([]model.ItemExtended, error)
	GetItemsBySellerID(ctx context.Context, sellerID int64) ([]model.Item, error)
	Delete(ctx context.Context, itemID int64) error
	Update(ctx context.Context, item model.Item) error
}

type Usecase struct {
	item itemRepo
}

func New(r itemRepo) *Usecase {
	return &Usecase{
		item: r,
	}
}

func (u *Usecase) Create(ctx context.Context, item model.Item) (int64, error) {
	err := model.ValidateItem(item)
	if err != nil {
		return 0, fmt.Errorf("item usecase: failed to validate item: %w", err)
	}

	itemID, err := u.item.Create(ctx, item)
	if err != nil {
		return 0, fmt.Errorf("item usecase: failed to create item: %w", err)
	}

	return itemID, nil
}

func (u *Usecase) GetByID(ctx context.Context, userID int64, itemID int64) (model.ItemExtended, error) {
	item, err := u.item.GetByID(ctx, itemID, userID)
	if err != nil {
		return model.ItemExtended{}, fmt.Errorf("item usecase: failed to get item: %w", err)
	}

	return item, nil
}

func (u *Usecase) GetAllItems(ctx context.Context, userID int64) ([]model.ItemExtended, error) {
	items, err := u.item.GetAllItems(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("item usecase: failed to get all items: %w", err)
	}

	return items, nil
}

func (u *Usecase) GetItemsBySellerID(ctx context.Context, sellerID int64) ([]model.Item, error) {
	items, err := u.item.GetItemsBySellerID(ctx, sellerID)
	if err != nil {
		return nil, fmt.Errorf("item usecase: failed to get seller items: %w", err)
	}

	return items, nil
}

func (u *Usecase) Delete(ctx context.Context, itemID int64) error {
	err := u.item.DeleteByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("item usecase: failed to delete item: %w", err)
	}

	return nil
}

func (u *Usecase) Update(ctx context.Context, item model.Item) error {
	err := model.ValidateItem(item)
	if err != nil {
		return fmt.Errorf("item usecase: invalid item: %w", err)
	}

	prevItem, err := u.item.GetByID(ctx, item.ID, item.Seller.ID)
	if err != nil {
		return fmt.Errorf("item usecase: failed to get item: %w", err)
	}

	if prevItem.Seller.ID != item.Seller.ID {
		return model.ErrNoAccess
	}

	err = u.item.UpdateByID(ctx, item)
	if err != nil {
		return fmt.Errorf("item usecase: failed to update item: %w", err)
	}

	return nil
}
