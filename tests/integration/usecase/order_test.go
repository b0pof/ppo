//go:build integration

package usecase

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	cartRepo "github.com/b0pof/ppo/internal/repository/cart"
	itemRepo "github.com/b0pof/ppo/internal/repository/item"
	orderRepo "github.com/b0pof/ppo/internal/repository/order"
	orderUsecase "github.com/b0pof/ppo/internal/usecase/order"
	"github.com/b0pof/ppo/tests/controller"
)

type OrderFlow struct {
	suite.Suite
}

func (g *OrderFlow) TestFullOrderFlow(t provider.T) {
	tests := []struct {
		name   string
		userID int64
	}{
		{
			name:   "full order usecase flow",
			userID: 1,
		},
	}

	for _, test := range tests {
		tt := test
		t.WithNewStep(tt.name, func(ctxA provider.StepCtx) {
			ctx := context.Background()
			ctrl := controller.NewController(t)

			orderRepository := orderRepo.New(ctrl.GetDB())
			itemRepository := itemRepo.New(ctrl.GetDB())
			cartRepository := cartRepo.New(ctrl.GetDB())

			usecase := orderUsecase.New(orderRepository, itemRepository, cartRepository)

			err := cartRepository.DeleteAllCartItems(ctx, tt.userID)
			ctxA.Assert().NoError(err)

			_, err = cartRepository.AddItemToCart(ctx, tt.userID, 1)
			ctxA.Assert().NoError(err)
			_, err = cartRepository.AddItemToCart(ctx, tt.userID, 3)
			ctxA.Assert().NoError(err)

			orderID, err := usecase.Create(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().True(orderID > 0)

			cartContent, err := cartRepository.GetCartContentByUserID(ctx, tt.userID)
			ctxA.Assert().ErrorIs(err, model.ErrNotFound)
			ctxA.Assert().Len(cartContent.Items, 0)

			orderExt, err := usecase.GetByID(ctx, orderID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(orderExt.Order.ID, orderID)
			ctxA.Assert().True(orderExt.Sum > 0)
			ctxA.Assert().Len(orderExt.Items, 2)

			allOrders, err := usecase.GetAllOrders(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			found := false
			for _, o := range allOrders {
				if o.ID == orderID {
					found = true
					break
				}
			}
			ctxA.Assert().True(found)

			newStatus := model.OrderStatusCreated
			err = usecase.UpdateStatus(ctx, orderID, newStatus)
			ctxA.Assert().NoError(err)

			orderAfterStatus, err := usecase.GetByID(ctx, orderID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(newStatus, orderAfterStatus.Status)

			err = usecase.Cancel(ctx, orderID)
			ctxA.Assert().NoError(err)

			orderAfterCancel, err := usecase.GetByID(ctx, orderID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(model.OrderStatusCancelled, orderAfterCancel.Status)
		})
	}
}
