package auth_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	. "github.com/b0pof/ppo/internal/usecase/auth"
)

type AuthUsecase struct {
	suite.Suite
}

func TestAuthSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthUsecase))
}

func (s *AuthUsecase) TestUsecase_Login(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		login        string
		password     string
		prepare      func(auth *MockauthRepo, user *MockuserRepo)
		expectations func(assert provider.Asserts, got string, err error)
	}{
		{
			name:     "success",
			login:    "test111",
			password: "testtest",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				user.EXPECT().GetByLogin(gomock.Any(), "test111").Return(model.User{
					ID:       1,
					Login:    "test111",
					Name:     "test",
					Role:     model.RoleBuyer,
					Phone:    "8-800-555-35-35",
					Password: "$2a$10$NbaNqa7I.TjwMdIp3z/uXeF4.8al6QWdX4CoFxEeQCN4R7vuQo7JW",
				}, nil)
				auth.EXPECT().CreateSession(int64(1)).Return("sess_hash")
			},
			expectations: func(assert provider.Asserts, got string, err error) {
				assert.NotEqual("", got)
				assert.NoError(err)
			},
		},
		{
			name:     "empty args",
			login:    "",
			password: "",
			expectations: func(assert provider.Asserts, got string, err error) {
				assert.ErrorIs(err, model.ErrInvalidInput)
			},
		},
		{
			name:     "failed to get user by login",
			login:    "test111",
			password: "testtest",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				user.EXPECT().GetByLogin(gomock.Any(), "test111").Return(model.User{}, errors.New("error"))
			},
			expectations: func(assert provider.Asserts, got string, err error) {
				assert.Error(err)
			},
		},
		{
			name:     "wrong password",
			login:    "test111",
			password: "testtest",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				user.EXPECT().GetByLogin(gomock.Any(), "test111").Return(model.User{
					ID:       1,
					Login:    "test111",
					Name:     "test",
					Role:     model.RoleBuyer,
					Phone:    "8-800-555-35-35",
					Password: ".2a$10$NbaNqa7I.TjwMdIp3z/uXeF4.8al6QWdX4CoFxEeQCN4R7vuQo7JW",
				}, nil)
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
			mockAuthRepo := NewMockauthRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockAuthRepo, mockUserRepo)
			}

			instance := New(mockAuthRepo, mockUserRepo)

			out, err := instance.Login(context.Background(), tc.login, tc.password)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *AuthUsecase) TestUsecase_Signup(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		login        string
		password     string
		prepare      func(auth *MockauthRepo, user *MockuserRepo)
		expectations func(assert provider.Asserts, gotSessionID string, err error)
	}{
		{
			name:     "success",
			login:    "test111",
			password: "testtest",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				user.EXPECT().Create(gomock.Any(), "test111", gomock.Any(), "buyer").Return(int64(1), nil)
				auth.EXPECT().CreateSession(int64(1)).Return("sess_hash")
			},
			expectations: func(assert provider.Asserts, gotSessionID string, err error) {
				assert.Equal("sess_hash", gotSessionID)
				assert.NoError(err)
			},
		},
		{
			name:     "empty args",
			login:    "",
			password: "",
			expectations: func(assert provider.Asserts, gotSessionID string, err error) {
				assert.ErrorIs(err, model.ErrInvalidInput)
			},
		},
		{
			name:     "failed to hash (password len > 72)",
			login:    "test111",
			password: strings.Repeat("abcd4", 15),
			expectations: func(assert provider.Asserts, gotSessionID string, err error) {
				assert.ErrorIs(err, model.ErrFailedToHash)
			},
		},
		{
			name:     "failed to create user",
			login:    "test111",
			password: "testtest",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				user.EXPECT().Create(gomock.Any(), "test111", gomock.Any(), "buyer").
					Return(int64(0), errors.New("error"))
			},
			expectations: func(assert provider.Asserts, gotSessionID string, err error) {
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockAuthRepo := NewMockauthRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockAuthRepo, mockUserRepo)
			}

			instance := New(mockAuthRepo, mockUserRepo)

			out, err := instance.Signup(context.Background(), tc.login, tc.password, "buyer")

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}

func (s *AuthUsecase) TestUsecase_Logout(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		sessionID    string
		prepare      func(auth *MockauthRepo, user *MockuserRepo)
		expectations func(assert provider.Asserts, err error)
	}{
		{
			name:      "success",
			sessionID: "sess_hash",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				auth.EXPECT().DeleteSession("sess_hash").Return(nil)
			},
			expectations: func(assert provider.Asserts, err error) {
				assert.NoError(err)
			},
		},
		{
			name:      "failed to delete session",
			sessionID: "sess_hash",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				auth.EXPECT().DeleteSession("sess_hash").Return(errors.New("error"))
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
			mockAuthRepo := NewMockauthRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockAuthRepo, mockUserRepo)
			}

			instance := New(mockAuthRepo, mockUserRepo)

			err := instance.Logout(tc.sessionID)

			tc.expectations(ctx.Assert(), err)
		})
	}
}

func (s *AuthUsecase) TestUsecase_IsLoggedIn(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		sessionID    string
		prepare      func(auth *MockauthRepo, user *MockuserRepo)
		expectations func(assert provider.Asserts, got bool)
	}{
		{
			name:      "success",
			sessionID: "sess_hash",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				auth.EXPECT().SessionExists("sess_hash").Return(true)
			},
			expectations: func(assert provider.Asserts, got bool) {
				assert.True(got)
			},
		},
		{
			name:      "no session found",
			sessionID: "sess_hash",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				auth.EXPECT().SessionExists("sess_hash").Return(false)
			},
			expectations: func(assert provider.Asserts, got bool) {
				assert.False(got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
			ctrl := gomock.NewController(t)
			mockAuthRepo := NewMockauthRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockAuthRepo, mockUserRepo)
			}

			instance := New(mockAuthRepo, mockUserRepo)

			out := instance.IsLoggedIn(tc.sessionID)

			tc.expectations(ctx.Assert(), out)
		})
	}
}

func (s *AuthUsecase) TestUsecase_GetUserIDBySessionID(t provider.T) {
	t.Parallel()

	tests := []struct {
		name         string
		sessionID    string
		prepare      func(auth *MockauthRepo, user *MockuserRepo)
		expectations func(assert provider.Asserts, got int64, err error)
	}{
		{
			name:      "success",
			sessionID: "sess_hash",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				auth.EXPECT().GetUserIDBySessionID("sess_hash").Return(int64(1), nil)
			},
			expectations: func(assert provider.Asserts, got int64, err error) {
				assert.NoError(err)
			},
		},
		{
			name:      "failed to get user by session",
			sessionID: "sess_hash",
			prepare: func(auth *MockauthRepo, user *MockuserRepo) {
				auth.EXPECT().GetUserIDBySessionID("sess_hash").
					Return(int64(0), errors.New("error"))
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
			mockAuthRepo := NewMockauthRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockAuthRepo, mockUserRepo)
			}

			instance := New(mockAuthRepo, mockUserRepo)

			out, err := instance.GetUserIDBySessionID(tc.sessionID)

			tc.expectations(ctx.Assert(), out, err)
		})
	}
}
