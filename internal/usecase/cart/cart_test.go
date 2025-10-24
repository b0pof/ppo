package cart_test

import (
	"context"
	"errors"
	"testing"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	. "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/cart"

	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type CartUsecase struct {
	suite.Suite
}

func TestCartSuite(t *testing.T) {
	suite.RunSuite(t, new(CartUsecase))
}

func (s *CartUsecase) TestUsecase_GetCartItemsAmount(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(cart *MockcartRepo)
		expectations func(assert provider.Asserts, got int, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().GetCartItemsAmount(gomock.Any(), int64(1)).Return(3, nil)
			},
			expectations: func(assert provider.Asserts, got int, err error) {
				assert.Equal(3, got)
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get cart item amoount",
			userID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().GetCartItemsAmount(gomock.Any(), int64(1)).Return(0, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got int, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockCartRepo)
			}

			instance := New(mockCartRepo)

			out, err := instance.GetCartItemsAmount(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *CartUsecase) TestUsecase_AddItem(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		itemID       int64
		prepare      func(cart *MockcartRepo)
		expectations func(assert provider.Asserts, got int, err error)
	}{
		{
			name:   "success",
			userID: 1,
			itemID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().AddItemToCart(gomock.Any(), int64(1), int64(1)).Return(1, nil)
			},
			expectations: func(assert provider.Asserts, got int, err error) {
				assert.Equal(1, got)
				assert.NoError(err)
			},
		},
		{
			name:   "failed to add item",
			userID: 1,
			itemID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().AddItemToCart(gomock.Any(), int64(1), int64(1)).Return(0, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got int, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockCartRepo)
			}

			instance := New(mockCartRepo)

			out, err := instance.AddItem(context.Background(), tc.userID, tc.itemID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *CartUsecase) TestUsecase_DeleteItem(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		itemID       int64
		prepare      func(cart *MockcartRepo)
		expectations func(assert provider.Asserts, got int, err error)
	}{
		{
			name:   "success",
			userID: 1,
			itemID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().DeleteCartItem(gomock.Any(), int64(1), int64(1)).Return(1, nil)
			},
			expectations: func(assert provider.Asserts, got int, err error) {
				assert.Equal(1, got)
				assert.NoError(err)
			},
		},
		{
			name:   "failed to delete item",
			userID: 1,
			itemID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().DeleteCartItem(gomock.Any(), int64(1), int64(1)).Return(0, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got int, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockCartRepo)
			}

			instance := New(mockCartRepo)

			out, err := instance.DeleteItem(context.Background(), tc.userID, tc.itemID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *CartUsecase) TestUsecase_Clear(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(cart *MockcartRepo)
		expectations func(assert provider.Asserts, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().DeleteAllCartItems(gomock.Any(), int64(1)).Return(nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.NoError(err)
			},
		},
		{
			name:   "failed to clear cart",
			userID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().DeleteAllCartItems(gomock.Any(), int64(1)).Return(errors.New("error"))
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockCartRepo)
			}

			instance := New(mockCartRepo)

			err := instance.Clear(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), err)
		})
	}
}

func (s *CartUsecase) TestUsecase_GetCartContent(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(cart *MockcartRepo)
		expectations func(assert provider.Asserts, got model.CartContent, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().GetCartContentByUserID(gomock.Any(), int64(1)).Return(model.CartContent{
					TotalPrice: 1500,
					TotalCount: 3,
					Items: []model.CartItem{
						{
							ID: 1,
						},
						{
							ID: 2,
						},
					},
				}, nil)
			},
			expectations: func(assert provider.Asserts, got model.CartContent, err error) {
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get cart content",
			userID: 1,
			prepare: func(cart *MockcartRepo) {
				cart.EXPECT().GetCartContentByUserID(gomock.Any(), int64(1)).
					Return(model.CartContent{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got model.CartContent, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockCartRepo)
			}

			instance := New(mockCartRepo)

			out, err := instance.GetCartContent(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}
