package repository_test

import (
	"context"
	"fmt"

	cartUc "github.com/b0pof/ppo/internal/usecase/cart"
	"github.com/b0pof/ppo/tests/integration/common/creator/user"
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
	t.Run("cart flow", func(t provider.T) {
		ctrl := controller.NewController(t)

		ctx := context.Background()
		testUser := user.New(ctrl).Random()

		itemRepository := itemRepo.New(ctrl.GetDB())
		cartRepository := cartRepo.New(ctrl.GetDB())
		cartUsecase := cartUc.New(cartRepository)

		testItem1 := model.Item{
			Name:        "Doll",
			Price:       1000,
			Description: "Great Doll!",
			ImgSrc:      "https://example.com/doll",
		}
		itemID1, err := itemRepository.Create(ctx, testItem1)
		t.Assert().NoError(err)

		testItem2 := model.Item{
			Name:        "Plastic car",
			Price:       1200,
			Description: "Nice car",
			ImgSrc:      "https://example.com/car",
		}
		itemID2, err := itemRepository.Create(ctx, testItem2)
		t.Assert().NoError(err)

		newCount1, err := cartUsecase.AddItem(ctx, testUser.ID, itemID1)
		t.Assert().NoError(err)

		newCount2, err := cartUsecase.AddItem(ctx, testUser.ID, itemID2)
		t.Assert().NoError(err)

		_, err = itemRepository.GetByID(ctx, testUser.ID, itemID1)
		t.Assert().NoError(err, fmt.Sprintf("NOT FOUND ITEM with id = %d", itemID1))

		cartItem2, err := itemRepository.GetByID(ctx, testUser.ID, itemID2)
		t.Assert().NoError(err)

		item1, err := itemRepository.GetByID(ctx, itemID1, testUser.ID)
		t.Assert().NoError(err)

		t.Assert().Equal(model.ItemExtended{
			Item: model.Item{
				ID:          itemID1,
				Name:        "Doll",
				Price:       1000,
				Description: "Great Doll!",
				ImgSrc:      "https://example.com/doll",
			},
			Amount: item1.Amount,
		}, item1)

		itemID3, err := itemRepository.Create(ctx, model.Item{
			Name:        "Test Item",
			Price:       500,
			Description: "Test",
			ImgSrc:      "https://example.com/test",
		})
		t.Assert().NoError(err)

		err = itemRepository.DeleteByID(ctx, itemID3)
		t.Assert().NoError(err)

		items, err := itemRepository.GetAllItems(ctx, testUser.ID)
		t.Assert().NoError(err)

		foundItems := 0
		for _, item := range items {
			if item.ID == itemID1 {
				foundItems++
			}
		}
		t.Assert().Equal(1, foundItems)

		cartItems, err := cartRepository.GetCartContentByUserID(ctx, testUser.ID)
		t.Assert().NoError(err)

		fmt.Println("cART ITEMS:", cartItems.Items)

		t.Assert().Contains(cartItems.Items, model.CartItem{
			ID:     itemID1,
			Name:   testItem1.Name,
			Price:  testItem1.Price,
			Count:  newCount1,
			ImgSrc: testItem1.ImgSrc,
		})

		t.Assert().Contains(cartItems.Items, model.CartItem{
			ID:     itemID2,
			Name:   cartItem2.Name,
			Price:  cartItem2.Price,
			Count:  newCount2,
			ImgSrc: cartItem2.ImgSrc,
		})
	})
}
