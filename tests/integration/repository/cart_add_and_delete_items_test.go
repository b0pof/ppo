//go:build integration

package repository

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	cartRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/cart"
	"git.iu7.bmstu.ru/kia22u475/ppo/tests/controller"
)

type CartFlow struct {
	suite.Suite
}

func (g *CartFlow) TestAddItemsToCart(t provider.T) {
	tests := []struct {
		name       string
		userID     int64
		itemsToBuy []model.CartItem
		expected   func(ctrl *controller.Controller, actual int, err error)
	}{
		{
			name:   "cart test",
			userID: 1,
			itemsToBuy: []model.CartItem{
				{
					ID:    1,
					Count: 1,
				},
				{
					ID:    3,
					Count: 2,
				},
			},
		},
	}

	for _, test := range tests {
		tt := test
		t.WithNewStep(tt.name, func(ctxA provider.StepCtx) {
			ctrl := controller.NewController(t)

			cartRepository := cartRepo.New(ctrl.GetDB())

			ctx := context.Background()

			// clear cart
			err := cartRepository.DeleteAllCartItems(ctx, tt.userID)
			ctxA.Assert().NoError(err)

			// check whether the cart is empty
			amount, err := cartRepository.GetCartItemsAmount(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(amount, 0)

			// add items to cart
			for _, item := range tt.itemsToBuy {

				for _ = range item.Count {
					amount, err = cartRepository.AddItemToCart(ctx, tt.userID, item.ID)
					ctxA.Assert().NoError(err)
				}

				// check the amount after the last addition
				ctxA.Assert().Equal(item.Count, amount)
			}

			// get cart items amount
			cartSize, err := cartRepository.GetCartItemsAmount(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(len(tt.itemsToBuy), cartSize)

			// decrease cart item amount
			newCount, err := cartRepository.DeleteCartItem(ctx, tt.userID, tt.itemsToBuy[0].ID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(tt.itemsToBuy[0].Count-1, newCount)

			// add back
			newCount, err = cartRepository.AddItemToCart(ctx, tt.userID, tt.itemsToBuy[0].ID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(tt.itemsToBuy[0].Count, newCount)
		})
	}
}
