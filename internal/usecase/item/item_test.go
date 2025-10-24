package item_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	. "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/item"
)

type ItemUsecase struct {
	suite.Suite
}

func TestItemSuite(t *testing.T) {
	suite.RunSuite(t, new(ItemUsecase))
}

func (s *ItemUsecase) TestUsecase_Create(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		item         model.Item
		prepare      func(item *MockitemRepo)
		expectations func(assert provider.Asserts, got int64, err error)
	}{
		{
			name: "success",
			item: model.Item{
				Name:        "test",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {
				item.EXPECT().Create(gomock.Any(), model.Item{
					Name:        "test",
					Description: "test",
					ImgSrc:      "https://test.com/img.png",
				}).Return(int64(1), nil)
			},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Equal(int64(1), got)
				assert.NoError(err)
			},
		},
		{
			name: "validation error: name length",
			item: model.Item{
				Name:        "tes",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: name pattern",
			item: model.Item{
				Name:        "test* test",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: description length",
			item: model.Item{
				Name:        "test",
				Description: "tes",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: description pattern",
			item: model.Item{
				Name:        "test",
				Description: "test* test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: imgSrc pattern",
			item: model.Item{
				Name:        "test",
				Description: "test",
				ImgSrc:      "not a link",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: imgSrc length",
			item: model.Item{
				Name:        "test",
				Description: "test",
				ImgSrc:      "not",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.Error(err)
			},
		},
		{
			name: "failed to create",
			item: model.Item{
				Name:        "test",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {
				item.EXPECT().Create(gomock.Any(), model.Item{
					Name:        "test",
					Description: "test",
					ImgSrc:      "https://test.com/img.png",
				}).Return(int64(0), errors.New("error"))
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
			mockItemRepo := NewMockitemRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockItemRepo)
			}

			instance := New(mockItemRepo)

			out, err := instance.Create(context.Background(), tc.item)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *ItemUsecase) TestUsecase_GetByID(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		itemID       int64
		prepare      func(item *MockitemRepo)
		expectations func(assert provider.Asserts, got model.ItemExtended, err error)
	}{
		{
			name:   "success",
			userID: 1,
			itemID: 1,
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetByID(gomock.Any(), int64(1), int64(1)).Return(model.ItemExtended{
					Item: model.Item{
						ID:          1,
						Name:        "test",
						Description: "test",
						ImgSrc:      "https://test.com/img.png",
						Seller: model.Seller{
							ID: 1,
						},
					},
				}, nil)
			},
			expectations: func(assert provider.Asserts, got model.ItemExtended, err error) {
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get by id",
			userID: 1,
			itemID: 1,
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetByID(gomock.Any(), int64(1), int64(1)).Return(model.ItemExtended{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got model.ItemExtended, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockItemRepo := NewMockitemRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockItemRepo)
			}

			instance := New(mockItemRepo)

			out, err := instance.GetByID(context.Background(), tc.userID, tc.itemID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *ItemUsecase) TestUsecase_GetAllItems(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(item *MockitemRepo)
		expectations func(assert provider.Asserts, got []model.ItemExtended, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetAllItems(gomock.Any(), int64(1)).Return([]model.ItemExtended{
					{
						Item: model.Item{
							ID:          1,
							Name:        "test",
							Description: "test",
							ImgSrc:      "https://test.com/img.png",
							Seller: model.Seller{
								ID: 1,
							},
						},
					},
					{
						Item: model.Item{
							ID:          2,
							Name:        "test2",
							Description: "test2",
							ImgSrc:      "https://test.com/img2.png",
							Seller: model.Seller{
								ID: 2,
							},
						},
					},
				}, nil)
			},
			expectations: func(assert provider.Asserts, got []model.ItemExtended, err error) {
				assert.Len(got, 2)
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get items",
			userID: 1,
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetAllItems(gomock.Any(), int64(1)).Return(nil, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got []model.ItemExtended, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockItemRepo := NewMockitemRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockItemRepo)
			}

			instance := New(mockItemRepo)

			out, err := instance.GetAllItems(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *ItemUsecase) TestUsecase_GetItemsBySellerID(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		sellerID     int64
		prepare      func(item *MockitemRepo)
		expectations func(assert provider.Asserts, got []model.Item, err error)
	}{
		{
			name:     "success",
			sellerID: 1,
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetItemsBySellerID(gomock.Any(), int64(1)).Return([]model.Item{
					{
						ID:          1,
						Name:        "test",
						Description: "test",
						ImgSrc:      "https://test.com/img.png",
						Seller: model.Seller{
							ID: 1,
						},
					},
					{
						ID:          2,
						Name:        "test2",
						Description: "test2",
						ImgSrc:      "https://test.com/img2.png",
						Seller: model.Seller{
							ID: 2,
						},
					},
				}, nil)
			},
			expectations: func(assert provider.Asserts, got []model.Item, err error) {
				assert.Len(got, 2)
				assert.NoError(err)
			},
		},
		{
			name:     "failed to get seller items",
			sellerID: 1,
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetItemsBySellerID(gomock.Any(), int64(1)).Return(nil, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got []model.Item, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockItemRepo := NewMockitemRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockItemRepo)
			}

			instance := New(mockItemRepo)

			out, err := instance.GetItemsBySellerID(context.Background(), tc.sellerID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *ItemUsecase) TestUsecase_Delete(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		itemID       int64
		prepare      func(item *MockitemRepo)
		expectations func(assert provider.Asserts, err error)
	}{
		{
			name:   "success",
			itemID: 1,
			prepare: func(item *MockitemRepo) {
				item.EXPECT().DeleteByID(gomock.Any(), int64(1)).Return(nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.NoError(err)
			},
		},
		{
			name:   "failed to delete",
			itemID: 1,
			prepare: func(item *MockitemRepo) {
				item.EXPECT().DeleteByID(gomock.Any(), int64(1)).Return(errors.New("error"))
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
			mockItemRepo := NewMockitemRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockItemRepo)
			}

			instance := New(mockItemRepo)

			err := instance.Delete(context.Background(), tc.itemID)

			tc.expectations(ctx.Assert(), err)
		})
	}
}

func (s *ItemUsecase) TestUsecase_Update(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		item         model.Item
		prepare      func(item *MockitemRepo)
		expectations func(assert provider.Asserts, err error)
	}{
		{
			name: "success",
			item: model.Item{
				Name:        "test",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetByID(gomock.Any(), int64(0), int64(0)).Return(model.ItemExtended{
					Item: model.Item{
						Name:        "test",
						Description: "test",
						ImgSrc:      "https://test.com/img.png",
					},
				}, nil)
				item.EXPECT().UpdateByID(gomock.Any(), model.Item{
					Name:        "test",
					Description: "test",
					ImgSrc:      "https://test.com/img.png",
				}).Return(nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.NoError(err)
			},
		},
		{
			name: "sellerID differs",
			item: model.Item{
				ID: 1,
				Seller: model.Seller{
					ID: 1,
				},
				Name:        "test",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetByID(gomock.Any(), int64(1), int64(1)).Return(model.ItemExtended{
					Item: model.Item{
						Seller: model.Seller{
							ID: 222,
						},
						Name:        "test",
						Description: "test",
						ImgSrc:      "https://test.com/img.png",
					},
				}, nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.ErrorIs(err, model.ErrNoAccess)
			},
		},
		{
			name: "failed to get item",
			item: model.Item{
				ID: 1,
				Seller: model.Seller{
					ID: 1,
				},
				Name:        "test",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetByID(gomock.Any(), int64(1), int64(1)).
					Return(model.ItemExtended{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: name length",
			item: model.Item{
				Name:        "tes",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: name pattern",
			item: model.Item{
				Name:        "test* test",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: description length",
			item: model.Item{
				Name:        "test",
				Description: "tes",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: description pattern",
			item: model.Item{
				Name:        "test",
				Description: "test* test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: imgSrc pattern",
			item: model.Item{
				Name:        "test",
				Description: "test",
				ImgSrc:      "not a link",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name: "validation error: imgSrc length",
			item: model.Item{
				Name:        "test",
				Description: "test",
				ImgSrc:      "not",
			},
			prepare: func(item *MockitemRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name: "failed to update",
			item: model.Item{
				Name:        "test",
				Description: "test",
				ImgSrc:      "https://test.com/img.png",
			},
			prepare: func(item *MockitemRepo) {
				item.EXPECT().GetByID(gomock.Any(), int64(0), int64(0)).Return(model.ItemExtended{
					Item: model.Item{
						Name:        "test",
						Description: "test",
						ImgSrc:      "https://test.com/img.png",
					},
				}, nil)
				item.EXPECT().UpdateByID(gomock.Any(), model.Item{
					Name:        "test",
					Description: "test",
					ImgSrc:      "https://test.com/img.png",
				}).Return(errors.New("error"))
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
			mockItemRepo := NewMockitemRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockItemRepo)
			}

			instance := New(mockItemRepo)

			err := instance.Update(context.Background(), tc.item)

			tc.expectations(ctx.Assert(), err)
		})
	}
}
