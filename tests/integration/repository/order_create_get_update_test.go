package repository_test

import (
	"context"
	"time"

	"github.com/b0pof/ppo/internal/model"
	orderRepo "github.com/b0pof/ppo/internal/repository/order"
	"github.com/b0pof/ppo/tests/controller"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

const (
	timeParseLayout = "2006-01-02 15:04:05.999999 -07:00"
)

type RepoOrderFlow struct {
	suite.Suite
}

func (g *RepoOrderFlow) TestOrder(t provider.T) {
	tests := []struct {
		name    string
		userID  int64
		orderID int64
	}{
		{
			name:    "order test",
			userID:  1,
			orderID: 1,
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t provider.T) {
			ctrl := controller.NewController(t)

			orderRepository := orderRepo.New(ctrl.GetDB())

			ctx := context.Background()

			// get order
			order, err := orderRepository.GetByID(ctx, tt.orderID)
			t.Assert().NoError(err)

			expectedTime, _ := time.Parse(timeParseLayout, "2025-03-23 10:04:24.284191 +00:00")
			location, _ := time.LoadLocation("Europe/Moscow")
			expectedTimeMsk := expectedTime.In(location).Local()

			t.Assert().Equal(model.Order{
				ID:         tt.orderID,
				Sum:        5000,
				BuyerID:    1,
				ItemsCount: 5,
				Status:     model.OrderStatusCreated,
				CreatedAt:  expectedTimeMsk,
			}, order)

			// get orders by user
			orders, err := orderRepository.GetOrdersByUserID(ctx, tt.userID)
			t.Assert().NoError(err)
			t.Assert().Equal([]model.Order{
				{
					ID:         1,
					Sum:        5000,
					BuyerID:    tt.userID,
					ItemsCount: 5,
					Status:     model.OrderStatusCreated,
					CreatedAt:  expectedTimeMsk,
				},
				{
					ID:         2,
					Sum:        6697,
					BuyerID:    tt.userID,
					ItemsCount: 4,
					Status:     model.OrderStatusReady,
					CreatedAt:  expectedTimeMsk,
				},
			}, orders)

			// update status
			err = orderRepository.UpdateStatus(ctx, tt.orderID, model.OrderStatusDone)
			t.Assert().NoError(err)

			// check new status
			order, err = orderRepository.GetByID(ctx, tt.orderID)
			t.Assert().NoError(err)
			t.Assert().Equal(model.Order{
				ID:         tt.orderID,
				Sum:        5000,
				BuyerID:    1,
				ItemsCount: 5,
				Status:     model.OrderStatusDone,
				CreatedAt:  expectedTimeMsk,
			}, order)
		})

	}
}
