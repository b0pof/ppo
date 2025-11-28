package usecase_test

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	authRepo "github.com/b0pof/ppo/internal/repository/auth"
	userRepo "github.com/b0pof/ppo/internal/repository/user"
	authUsecase "github.com/b0pof/ppo/internal/usecase/auth"
	"github.com/b0pof/ppo/tests/controller"
)

type AuthTest struct {
	suite.Suite
}

func (g *AuthTest) TestFullAuthFlow(t provider.T) {
	tests := []struct {
		name     string
		login    string
		password string
		role     string
	}{
		{
			name:     "auth usecase full flow",
			login:    "test_user_login",
			password: "TestPass123!",
			role:     "buyer",
		},
	}

	for _, test := range tests {
		tt := test
		t.WithNewStep(tt.name, func(ctxA provider.StepCtx) {
			ctx := context.Background()
			ctrl := controller.NewController(t)

			authR := authRepo.New(ctrl.GetRedis())
			userR := userRepo.New(ctrl.GetDB())
			usecase := authUsecase.New(authR, userR)

			sessionID, err := usecase.Signup(ctx, tt.login, tt.password, tt.role)
			ctxA.Assert().NoError(err)
			ctxA.Assert().NotEmpty(sessionID)

			user, err := userR.GetByLogin(ctx, tt.login)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(tt.role, user.Role)

			ctxA.Assert().True(usecase.IsLoggedIn(sessionID))

			userID, err := usecase.GetUserIDBySessionID(sessionID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(user.ID, userID)

			err = usecase.Logout(sessionID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().False(usecase.IsLoggedIn(sessionID))

			newSessionID, err := usecase.Login(ctx, tt.login, tt.password)
			ctxA.Assert().NoError(err)
			ctxA.Assert().NotEmpty(newSessionID)
			ctxA.Assert().True(usecase.IsLoggedIn(newSessionID))

			_, err = usecase.Login(ctx, tt.login, "wrongpassword")
			ctxA.Assert().Error(err)
			ctxA.Assert().ErrorIs(err, model.ErrWrongPassword)

			_, err = usecase.Login(ctx, "", "")
			ctxA.Assert().Error(err)
			ctxA.Assert().ErrorIs(err, model.ErrInvalidInput)

			_, err = usecase.Signup(ctx, "", "", "")
			ctxA.Assert().Error(err)
			ctxA.Assert().ErrorIs(err, model.ErrInvalidInput)
		})
	}
}
