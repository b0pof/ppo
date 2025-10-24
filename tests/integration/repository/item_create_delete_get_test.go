//go:build integration

package repository

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	cartRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/cart"
	itemRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/item"
	"git.iu7.bmstu.ru/kia22u475/ppo/tests/controller"
)

type ItemFlow struct {
	suite.Suite
}

func (g *ItemFlow) TestGetCartItemsAmount(t provider.T) {
	tests := []struct {
		name   string
		userID int64
		itemID int64
	}{
		{
			name:   "test item repository",
			userID: 1,
			itemID: 1,
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t provider.T) {
			ctrl := controller.NewController(t)

			ctx := context.Background()

			itemRepository := itemRepo.New(ctrl.GetDB())
			cartRepository := cartRepo.New(ctrl.GetDB())

			// get item
			item, err := itemRepository.GetByID(ctx, tt.itemID, tt.userID)
			t.Assert().NoError(err)
			t.Assert().Equal(model.ItemExtended{
				Item: model.Item{
					ID:   tt.itemID,
					Name: "Doll",
					Seller: model.Seller{
						ID:   2,
						Name: "User2 Name",
					},
					Price:       1000,
					Description: "Great Doll!",
					ImgSrc:      "https://example.com/doll",
				},
				SellerName: "User2 Name",
				Amount:     3,
			}, item)

			// delete item
			err = itemRepository.DeleteByID(ctx, int64(4))
			t.Assert().NoError(err)

			// fetch all items
			items, err := itemRepository.GetAllItems(ctx, tt.userID)
			t.Assert().NoError(err)
			t.Assert().Len(items, 3)

			// get cart items
			cartItems, err := cartRepository.GetCartContentByUserID(ctx, tt.userID)
			t.Assert().NoError(err)
			t.Assert().Equal([]model.CartItem{
				{
					ID:     1,
					Name:   "Doll",
					Price:  1000,
					Count:  3,
					ImgSrc: "https://example.com/doll",
				},
				{
					ID:     2,
					Name:   "Plastic car",
					Price:  1200,
					Count:  1,
					ImgSrc: "https://example.com/car",
				},
			}, cartItems.Items)
		})
	}
}
