package repository_test

import (
	"context"
	"fmt"

	cartRepo "github.com/b0pof/ppo/internal/repository/cart"
	itemRepo "github.com/b0pof/ppo/internal/repository/item"
	cartUc "github.com/b0pof/ppo/internal/usecase/cart"
	itemUc "github.com/b0pof/ppo/internal/usecase/item"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	orderRepo "github.com/b0pof/ppo/internal/repository/order"
	"github.com/b0pof/ppo/tests/controller"
	"github.com/b0pof/ppo/tests/integration/common/creator/user"
)

const (
	timeParseLayout = "2006-01-02 15:04:05.999999 -07:00"
)

type RepoOrderFlow struct {
	suite.Suite
}

func (g *RepoOrderFlow) TestOrder(t provider.T) {
	t.Run("order test", func(t provider.T) {
		ctrl := controller.NewController(t)
		orderRepository := orderRepo.New(ctrl.GetDB())

		ctx := context.Background()
		testUser := user.New(ctrl).Random()

		cartUsecase := cartUc.New(cartRepo.New(ctrl.GetDB()))
		itemUsecase := itemUc.New(itemRepo.New(ctrl.GetDB()))

		var cartItems []model.CartItem

		for i := 0; i < 5; i++ {
			itemID, err := itemUsecase.Create(ctx, model.Item{
				Name:  fmt.Sprintf("item%d", i+1),
				Price: 1_000 * (i + 1),
			})
			t.Assert().NoError(err)

			newCount, err := cartUsecase.AddItem(ctx, testUser.ID, itemID)
			t.Assert().NoError(err)

			item, err := itemUsecase.GetByID(ctx, testUser.ID, itemID)
			t.Assert().NoError(err, fmt.Sprintf("item with id=%d not found", itemID))

			cartItems = append(cartItems, model.CartItem{
				ID:     item.ID,
				Name:   item.Name,
				Price:  item.Price,
				Count:  newCount,
				ImgSrc: item.ImgSrc,
			})
		}

		orderID, err := orderRepository.Create(ctx, testUser.ID, cartItems)
		t.Assert().NoError(err)

		order, err := orderRepository.GetByID(ctx, orderID)
		t.Assert().NoError(err)

		t.Assert().Equal(model.Order{
			ID:         orderID,
			Sum:        15_000,
			BuyerID:    testUser.ID,
			ItemsCount: 5,
			Status:     model.OrderStatusCreated,
			CreatedAt:  order.CreatedAt,
		}, order)

		orders, err := orderRepository.GetOrdersByUserID(ctx, testUser.ID)
		t.Assert().NoError(err)

		t.Assert().Equal([]model.Order{
			{
				ID:         orderID,
				Sum:        15_000,
				BuyerID:    testUser.ID,
				ItemsCount: 5,
				Status:     model.OrderStatusCreated,
				CreatedAt:  order.CreatedAt,
			},
		}, orders)

		err = orderRepository.UpdateStatus(ctx, orderID, model.OrderStatusDone)
		t.Assert().NoError(err)

		order, err = orderRepository.GetByID(ctx, orderID)
		t.Assert().NoError(err)

		t.Assert().Equal(model.Order{
			ID:         orderID,
			Sum:        15_000,
			BuyerID:    testUser.ID,
			ItemsCount: 5,
			Status:     model.OrderStatusDone,
			CreatedAt:  order.CreatedAt,
		}, order)
	})
}
