package order_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	. "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/order"

	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type OrderUsecase struct {
	suite.Suite
}

func TestOrderSuite(t *testing.T) {
	suite.RunSuite(t, new(OrderUsecase))
}

func (s *OrderUsecase) TestUsecase_Create(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(order *MockorderRepo, item *MockitemRepo, cart *MockcartRepo)
		expectations func(assert provider.Asserts, got int64, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo, cart *MockcartRepo) {
				cart.EXPECT().GetCartContentByUserID(gomock.Any(), int64(1)).Return(model.CartContent{
					TotalCount: 4,
					TotalPrice: 1500,
					Items: []model.CartItem{
						{
							ID:   1,
							Name: "first",
						},
						{
							ID:   2,
							Name: "second",
						},
					},
				}, nil)
				order.EXPECT().Create(gomock.Any(), int64(1), []model.CartItem{
					{
						ID:   1,
						Name: "first",
					},
					{
						ID:   2,
						Name: "second",
					},
				}).Return(int64(1), nil)
				cart.EXPECT().DeleteAllCartItems(gomock.Any(), int64(1)).Return(nil)
			},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Equal(int64(1), got)
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get cart items",
			userID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo, cart *MockcartRepo) {
				cart.EXPECT().GetCartContentByUserID(gomock.Any(), int64(1)).
					Return(model.CartContent{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
		{
			name:   "no items in cart",
			userID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo, cart *MockcartRepo) {
				cart.EXPECT().GetCartContentByUserID(gomock.Any(), int64(1)).
					Return(model.CartContent{}, nil)
			},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.ErrorIs(err, model.ErrCartIsEmpty)
			},
		},
		{
			name:   "failed to create order",
			userID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo, cart *MockcartRepo) {
				cart.EXPECT().GetCartContentByUserID(gomock.Any(), int64(1)).Return(model.CartContent{
					TotalCount: 4,
					TotalPrice: 1500,
					Items: []model.CartItem{
						{
							ID:   1,
							Name: "first",
						},
						{
							ID:   2,
							Name: "second",
						},
					},
				}, nil)
				order.EXPECT().Create(gomock.Any(), int64(1), []model.CartItem{
					{
						ID:   1,
						Name: "first",
					},
					{
						ID:   2,
						Name: "second",
					},
				}).Return(int64(0), errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
		{
			name:   "failed to clear cart",
			userID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo, cart *MockcartRepo) {
				cart.EXPECT().GetCartContentByUserID(gomock.Any(), int64(1)).Return(model.CartContent{
					TotalCount: 4,
					TotalPrice: 1500,
					Items: []model.CartItem{
						{
							ID:   1,
							Name: "first",
						},
						{
							ID:   2,
							Name: "second",
						},
					},
				}, nil)
				order.EXPECT().Create(gomock.Any(), int64(1), []model.CartItem{
					{
						ID:   1,
						Name: "first",
					},
					{
						ID:   2,
						Name: "second",
					},
				}).Return(int64(1), nil)
				cart.EXPECT().DeleteAllCartItems(gomock.Any(), int64(1)).Return(errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockOrderRepo := NewMockorderRepo(ctrl)
			mockItemRepo := NewMockitemRepo(ctrl)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockOrderRepo, mockItemRepo, mockCartRepo)
			}

			instance := New(mockOrderRepo, mockItemRepo, mockCartRepo)

			out, err := instance.Create(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *OrderUsecase) TestUsecase_GetByID(t provider.T) {
	t.Parallel()

	testTime := time.Now()

	tests := []struct {
		name         string
		orderID      int64
		prepare      func(order *MockorderRepo, item *MockitemRepo)
		expectations func(assert provider.Asserts, got model.OrderExtended, err error)
	}{
		{
			name:    "success",
			orderID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.Order{
					ID:         1,
					ItemsCount: 3,
					BuyerID:    2,
					Status:     model.OrderStatusCreated,
					CreatedAt:  testTime,
				}, nil)
				item.EXPECT().GetItemsByOrderID(gomock.Any(), int64(1)).Return([]model.OrderItemInfo{
					{
						ID:          111,
						ProductName: "test1",
						Price:       2000,
						Count:       1,
						ImgSrc:      "https://test.com/img1.png",
					},
					{
						ID:          222,
						ProductName: "test2",
						Price:       1000,
						Count:       3,
						ImgSrc:      "https://test.com/img2.png",
					},
				}, nil)
			},
			expectations: func(assert provider.Asserts, got model.OrderExtended, err error) {
				assert.Equal(model.OrderExtended{
					Order: model.Order{
						ID:         1,
						ItemsCount: 3,
						BuyerID:    2,
						Status:     model.OrderStatusCreated,
						CreatedAt:  testTime,
					},
					Sum: 5000,
					Items: []model.OrderItemInfo{
						{
							ID:          111,
							ProductName: "test1",
							Price:       2000,
							Count:       1,
							ImgSrc:      "https://test.com/img1.png",
						},
						{
							ID:          222,
							ProductName: "test2",
							Price:       1000,
							Count:       3,
							ImgSrc:      "https://test.com/img2.png",
						},
					},
				}, got)
				assert.NoError(err)
			},
		},
		{
			name:    "failed to get order",
			orderID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.Order{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got model.OrderExtended, err error) {
				assert.Error(err)
			},
		},
		{
			name:    "failed to get order items",
			orderID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.Order{
					ID:         1,
					ItemsCount: 3,
					BuyerID:    2,
					Status:     model.OrderStatusCreated,
					CreatedAt:  testTime,
				}, nil)
				item.EXPECT().GetItemsByOrderID(gomock.Any(), int64(1)).Return(nil, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got model.OrderExtended, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockOrderRepo := NewMockorderRepo(ctrl)
			mockItemRepo := NewMockitemRepo(ctrl)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockOrderRepo, mockItemRepo)
			}

			instance := New(mockOrderRepo, mockItemRepo, mockCartRepo)

			out, err := instance.GetByID(context.Background(), tc.orderID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *OrderUsecase) TestUsecase_GetAllOrders(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(order *MockorderRepo, item *MockitemRepo)
		expectations func(assert provider.Asserts, got []model.Order, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetOrdersByUserID(gomock.Any(), int64(1)).Return([]model.Order{
					{
						ID:         1,
						ItemsCount: 3,
						BuyerID:    2,
						Status:     model.OrderStatusCreated,
						CreatedAt:  time.Now(),
					},
				}, nil)
			},
			expectations: func(assert provider.Asserts, got []model.Order, err error) {
				assert.Len(got, 1)
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get orders",
			userID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetOrdersByUserID(gomock.Any(), int64(1)).Return(nil, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got []model.Order, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockOrderRepo := NewMockorderRepo(ctrl)
			mockItemRepo := NewMockitemRepo(ctrl)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockOrderRepo, mockItemRepo)
			}

			instance := New(mockOrderRepo, mockItemRepo, mockCartRepo)

			out, err := instance.GetAllOrders(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *OrderUsecase) TestUsecase_UpdateStatus(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		itemID       int64
		newStatus    string
		prepare      func(order *MockorderRepo, item *MockitemRepo)
		expectations func(assert provider.Asserts, err error)
	}{
		{
			name:      "success",
			itemID:    1,
			newStatus: model.OrderStatusReady,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.Order{
					ID: 1,
				}, nil)
				order.EXPECT().UpdateStatus(gomock.Any(), int64(1), model.OrderStatusReady).Return(nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.NoError(err)
			},
		},
		{
			name:      "failed to get order",
			itemID:    1,
			newStatus: model.OrderStatusReady,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).
					Return(model.Order{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:      "failed to update status",
			itemID:    1,
			newStatus: model.OrderStatusReady,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.Order{
					ID: 1,
				}, nil)
				order.EXPECT().UpdateStatus(gomock.Any(), int64(1), model.OrderStatusReady).
					Return(errors.New("error"))
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
			mockOrderRepo := NewMockorderRepo(ctrl)
			mockItemRepo := NewMockitemRepo(ctrl)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockOrderRepo, mockItemRepo)
			}

			instance := New(mockOrderRepo, mockItemRepo, mockCartRepo)

			err := instance.UpdateStatus(context.Background(), tc.itemID, tc.newStatus)

			tc.expectations(ctx.Assert(), err)
		})
	}
}

func (s *OrderUsecase) TestUsecase_Cancel(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		itemID       int64
		prepare      func(order *MockorderRepo, item *MockitemRepo)
		expectations func(assert provider.Asserts, err error)
	}{
		{
			name:   "success",
			itemID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.Order{
					ID: 1,
				}, nil)
				order.EXPECT().UpdateStatus(gomock.Any(), int64(1), model.OrderStatusCancelled).Return(nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get order",
			itemID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.Order{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:   "success",
			itemID: 1,
			prepare: func(order *MockorderRepo, item *MockitemRepo) {
				order.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.Order{
					ID: 1,
				}, nil)
				order.EXPECT().UpdateStatus(gomock.Any(), int64(1), model.OrderStatusCancelled).
					Return(errors.New("error"))
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
			mockOrderRepo := NewMockorderRepo(ctrl)
			mockItemRepo := NewMockitemRepo(ctrl)
			mockCartRepo := NewMockcartRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockOrderRepo, mockItemRepo)
			}

			instance := New(mockOrderRepo, mockItemRepo, mockCartRepo)

			err := instance.Cancel(context.Background(), tc.itemID)

			tc.expectations(ctx.Assert(), err)
		})
	}
}
