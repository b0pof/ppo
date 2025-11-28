package usecase_test

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	itemRepo "github.com/b0pof/ppo/internal/repository/item"
	itemUsecase "github.com/b0pof/ppo/internal/usecase/item"
	"github.com/b0pof/ppo/tests/controller"
)

type ItemFlow struct {
	suite.Suite
}

func (g *ItemFlow) TestFullItemFlow(t provider.T) {
	tests := []struct {
		name     string
		userID   int64
		sellerID int64
	}{
		{
			name:     "full item usecase flow",
			userID:   1,
			sellerID: 2,
		},
	}

	for _, test := range tests {
		tt := test
		t.WithNewStep(tt.name, func(ctxA provider.StepCtx) {
			ctx := context.Background()
			ctrl := controller.NewController(t)

			repo := itemRepo.New(ctrl.GetDB())
			usecase := itemUsecase.New(repo)

			itemToCreate := model.Item{
				Name:        "Toy Horse",
				Price:       2500,
				Description: "Nice toy horse for kids",
				ImgSrc:      "https://example.com/horse",
				Seller: model.Seller{
					ID:   tt.sellerID,
					Name: "User2 Name",
				},
			}

			itemID, err := usecase.Create(ctx, itemToCreate)
			ctxA.Assert().NoError(err)
			ctxA.Assert().True(itemID > 0)

			createdItem, err := usecase.GetByID(ctx, tt.userID, itemID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(itemToCreate.Name, createdItem.Name)
			ctxA.Assert().Equal(itemToCreate.Price, createdItem.Price)
			ctxA.Assert().Equal(itemToCreate.Seller.ID, createdItem.Seller.ID)

			allItems, err := usecase.GetAllItems(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			found := false
			for _, it := range allItems {
				if it.ID == itemID {
					found = true
					break
				}
			}
			ctxA.Assert().True(found)

			sellerItems, err := usecase.GetItemsBySellerID(ctx, tt.sellerID)
			ctxA.Assert().NoError(err)
			hasItem := false
			for _, it := range sellerItems {
				if it.ID == itemID {
					hasItem = true
					break
				}
			}
			ctxA.Assert().True(hasItem)

			updatedItem := model.Item{
				ID:          itemID,
				Name:        "Toy Horse Deluxe",
				Price:       3000,
				Description: "Upgraded toy horse",
				ImgSrc:      itemToCreate.ImgSrc,
				Seller: model.Seller{
					ID: tt.sellerID,
				},
			}
			err = usecase.Update(ctx, updatedItem)
			ctxA.Assert().NoError(err)

			afterUpdate, err := usecase.GetByID(ctx, tt.userID, itemID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(updatedItem.Name, afterUpdate.Name)
			ctxA.Assert().Equal(updatedItem.Price, afterUpdate.Price)
			ctxA.Assert().Equal(updatedItem.Description, afterUpdate.Description)

			wrongSellerItem := updatedItem
			wrongSellerItem.Seller.ID = 999 // другой селлер
			err = usecase.Update(ctx, wrongSellerItem)
			ctxA.Assert().Error(err)
			ctxA.Assert().ErrorIs(err, model.ErrNoAccess)

			err = usecase.Delete(ctx, itemID)
			ctxA.Assert().NoError(err)

			_, err = usecase.GetByID(ctx, tt.userID, itemID)
			ctxA.Assert().Error(err)
		})
	}
}
