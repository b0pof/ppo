package usecase_test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
	t.WithNewStep("auth flow", func(ctxA provider.StepCtx) {
		ctx := context.Background()
		ctrl := controller.NewController(t)
		// defer ctrl.Finish()

		authR := authRepo.New(ctrl.GetRedis())
		userR := userRepo.New(ctrl.GetDB())
		usecase := authUsecase.New(authR, userR)

		testUser := model.User{
			Login:    uuid.New().String(),
			Name:     uuid.New().String(),
			Role:     model.RoleBuyer,
			Password: uuid.New().String(),
		}

		sessionID, err := usecase.Signup(ctx, testUser.Login, testUser.Password, testUser.Role)
		ctxA.Assert().NoError(err)
		ctxA.Assert().NotEmpty(sessionID)

		user, err := userR.GetByLogin(ctx, testUser.Login)
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(testUser.Role, user.Role)

		ctxA.Assert().True(usecase.IsLoggedIn(sessionID))

		userID, err := usecase.GetUserIDBySessionID(sessionID)
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(user.ID, userID)

		err = usecase.Logout(sessionID)
		ctxA.Assert().NoError(err, fmt.Sprintf("NOT FOUND: session id = %s", sessionID))
		ctxA.Assert().False(usecase.IsLoggedIn(sessionID))

		newSessionID, err := usecase.Login(ctx, testUser.Login, testUser.Password)
		ctxA.Assert().NoError(err)
		ctxA.Assert().NotEmpty(newSessionID)
		ctxA.Assert().True(usecase.IsLoggedIn(newSessionID))

		_, err = usecase.Login(ctx, testUser.Login, "wrongpassword")
		ctxA.Assert().ErrorIs(err, model.ErrWrongPassword)

		_, err = usecase.Login(ctx, "", "")
		ctxA.Assert().Error(err)
		ctxA.Assert().ErrorIs(err, model.ErrInvalidInput)

		_, err = usecase.Signup(ctx, "", "", "")
		ctxA.Assert().Error(err)
		ctxA.Assert().ErrorIs(err, model.ErrInvalidInput)
	})
}
