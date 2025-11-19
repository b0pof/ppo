package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	. "github.com/b0pof/ppo/internal/usecase/user"
)

type UserUsecase struct {
	suite.Suite
}

func TestUserSuite(t *testing.T) {
	suite.RunSuite(t, new(UserUsecase))
}

func (s *UserUsecase) TestUsecase_GetByID(t provider.T) {
	t.Parallel()

	testUser := NewUserBuilder().
		WithRole(model.RoleBuyer).
		Build()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(user *MockuserRepo)
		expectations func(assert provider.Asserts, got model.User, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetByID(gomock.Any(), int64(1)).Return(testUser, nil)
			},
			expectations: func(assert provider.Asserts, got model.User, err error) {
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get user",
			userID: 1,
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.User{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got model.User, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockUserRepo)
			}

			instance := New(mockUserRepo)

			out, err := instance.GetByID(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *UserUsecase) TestUsecase_GetRoleByID(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(user *MockuserRepo)
		expectations func(assert provider.Asserts, got string, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetRoleByID(gomock.Any(), int64(1)).Return(model.RoleBuyer, nil)
			},
			expectations: func(assert provider.Asserts, got string, err error) {
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get role",
			userID: 1,
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetRoleByID(gomock.Any(), int64(1)).Return("", errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got string, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockUserRepo)
			}

			instance := New(mockUserRepo)

			out, err := instance.GetRoleByID(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *UserUsecase) TestUsecase_GetMetaInfoByUserID(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		prepare      func(user *MockuserRepo)
		expectations func(assert provider.Asserts, got model.UserMetaInfo, err error)
	}{
		{
			name:   "success",
			userID: 1,
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetUserMetaByID(gomock.Any(), int64(1)).Return(model.UserMetaInfo{
					Name:            "test",
					CartItemsAmount: 3,
				}, nil)
			},
			expectations: func(assert provider.Asserts, got model.UserMetaInfo, err error) {
				assert.NoError(err)
			},
		},
		{
			name:   "failed to get meta info",
			userID: 1,
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetUserMetaByID(gomock.Any(), int64(1)).Return(model.UserMetaInfo{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got model.UserMetaInfo, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockUserRepo)
			}

			instance := New(mockUserRepo)

			out, err := instance.GetMetaInfoByUserID(context.Background(), tc.userID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *UserUsecase) TestUsecase_UpdateByID(t provider.T) {
	t.Parallel()

	testUser := NewUserBuilder().
		WithRole(model.RoleBuyer).
		Build()

	tests := []struct {
		name         string
		userID       int64
		user         model.User
		prepare      func(user *MockuserRepo)
		expectations func(assert provider.Asserts, err error)
	}{
		{
			name:   "success",
			userID: 1,
			user:   testUser,
			prepare: func(user *MockuserRepo) {
				user.EXPECT().UpdateByID(gomock.Any(), int64(1), testUser).Return(nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.NoError(err)
			},
		},
		{
			name:    "validation: name length",
			userID:  1,
			user:    NewUserBuilder().WithName("").Build(),
			prepare: func(user *MockuserRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:    "validation: name pattern",
			userID:  1,
			user:    NewUserBuilder().WithName("name* name").Build(),
			prepare: func(user *MockuserRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:    "validation: login length",
			userID:  1,
			user:    NewUserBuilder().WithLogin("tes").Build(),
			prepare: func(user *MockuserRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:    "validation: login pattern",
			userID:  1,
			user:    NewUserBuilder().WithLogin("test* test").Build(),
			prepare: func(user *MockuserRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:    "validation: phone pattern",
			userID:  1,
			user:    NewUserBuilder().WithPhone("invalid").Build(),
			prepare: func(user *MockuserRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:   "failed to update",
			userID: 1,
			user:   testUser,
			prepare: func(user *MockuserRepo) {
				user.EXPECT().UpdateByID(gomock.Any(), int64(1), testUser).Return(errors.New("error"))
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
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockUserRepo)
			}

			instance := New(mockUserRepo)

			err := instance.UpdateByID(context.Background(), tc.userID, tc.user)

			tc.expectations(ctx.Assert(), err)
		})
	}
}

func (s *UserUsecase) TestUsecase_UpdatePassword(t provider.T) {
	t.Parallel()

	testUser := NewUserBuilder().
		WithPassword("$2a$10$NbaNqa7I.TjwMdIp3z/uXeF4.8al6QWdX4CoFxEeQCN4R7vuQo7JW").
		Build()

	tests := []struct {
		name         string
		userID       int64
		oldPassword  string
		password     string
		prepare      func(user *MockuserRepo)
		expectations func(assert provider.Asserts, err error)
	}{
		{
			name:        "success",
			userID:      1,
			oldPassword: "testtest",
			password:    "testtest22",
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetByID(gomock.Any(), int64(1)).
					Return(testUser, nil)
				user.EXPECT().UpdatePasswordByID(gomock.Any(), int64(1), gomock.Any()).Return(nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.NoError(err)
			},
		},
		{
			name:        "validation: password length",
			userID:      1,
			oldPassword: "testtest",
			password:    "test",
			prepare:     func(user *MockuserRepo) {},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:        "failed to get user",
			userID:      1,
			oldPassword: "testtest",
			password:    "testtest22",
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetByID(gomock.Any(), int64(1)).Return(model.User{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.Error(err)
			},
		},
		{
			name:        "wrong password",
			userID:      1,
			oldPassword: "testtest",
			password:    "testtest22",
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetByID(gomock.Any(), int64(1)).
					Return(NewUserBuilder().
						WithPassword(".2a$10$NbaNqa7I.TjwMdIp3z/uXeF4.8al6QWdX4CoFxEeQCN4R7vuQo7JW").
						Build(), nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.ErrorIs(err, model.ErrWrongPassword)
			},
		},
		{
			name:        "success",
			userID:      1,
			oldPassword: "testtest",
			password:    "testtest22",
			prepare: func(user *MockuserRepo) {
				user.EXPECT().GetByID(gomock.Any(), int64(1)).Return(testUser, nil)
				user.EXPECT().UpdatePasswordByID(gomock.Any(), int64(1), gomock.Any()).Return(errors.New("error"))
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
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockUserRepo)
			}

			instance := New(mockUserRepo)

			err := instance.UpdatePassword(context.Background(), tc.userID, tc.oldPassword, tc.password)

			tc.expectations(ctx.Assert(), err)
		})
	}
}
