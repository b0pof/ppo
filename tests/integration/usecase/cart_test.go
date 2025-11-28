package usecase_test

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	cartRepo "github.com/b0pof/ppo/internal/repository/cart"
	cartUsecase "github.com/b0pof/ppo/internal/usecase/cart"
	"github.com/b0pof/ppo/tests/controller"
)

type CartTest struct {
	suite.Suite
}

func (g *CartTest) TestFullCartFlow(t provider.T) {
	tests := []struct {
		name    string
		userID  int64
		itemID1 int64
		itemID2 int64
	}{
		{
			name:    "full cart usecase flow",
			userID:  1,
			itemID1: 2,
			itemID2: 3,
		},
	}

	for _, test := range tests {
		tt := test
		t.WithNewStep(tt.name, func(ctxA provider.StepCtx) {
			ctx := context.Background()
			ctrl := controller.NewController(t)

			repo := cartRepo.New(ctrl.GetDB())
			usecase := cartUsecase.New(repo)

			content, err := usecase.GetCartContent(ctx, tt.userID)
			ctxA.Assert().NoError(err)

			itemsCount := make(map[int64]int, len(content.Items))
			for _, item := range content.Items {
				itemsCount[item.ID] = item.Count
			}

			newAmount, err := usecase.AddItem(ctx, tt.userID, tt.itemID1)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(itemsCount[tt.itemID1]+1, newAmount)

			newAmount, err = usecase.AddItem(ctx, tt.userID, tt.itemID2)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(itemsCount[tt.itemID2]+1, newAmount)

			content, err = usecase.GetCartContent(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().True(len(content.Items) >= 2)

			newAmount, err = usecase.DeleteItem(ctx, tt.userID, tt.itemID1)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(itemsCount[tt.itemID1], newAmount)

			err = usecase.Clear(ctx, tt.userID)
			ctxA.Assert().NoError(err)

			finalAmount, err := usecase.GetCartItemsAmount(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(0, finalAmount)
		})
	}
}
