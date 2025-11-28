package repository_test

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	cartRepo "github.com/b0pof/ppo/internal/repository/cart"
	itemRepo "github.com/b0pof/ppo/internal/repository/item"
	"github.com/b0pof/ppo/tests/controller"
)

type RepoItemFlow struct {
	suite.Suite
}

func (g *RepoItemFlow) TestGetCartItemsAmount(t provider.T) {
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
			t.Assert().Contains(cartItems.Items, model.CartItem{
				ID:     1,
				Name:   "Doll",
				Price:  1000,
				Count:  3,
				ImgSrc: "https://example.com/doll",
			})
			t.Assert().Contains(cartItems.Items, model.CartItem{
				ID:     2,
				Name:   "Plastic car",
				Price:  1200,
				Count:  1,
				ImgSrc: "https://example.com/car",
			})
		})
	}
}
