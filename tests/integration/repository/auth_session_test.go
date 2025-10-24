//go:build integration

package repository

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	authRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/auth"
	"git.iu7.bmstu.ru/kia22u475/ppo/tests/controller"
)

type SessionCreateDelete struct {
	suite.Suite
}

func (g *SessionCreateDelete) TestLoginLogout(t provider.T) {
	tests := []struct {
		name     string
		userID   int64
		login    string
		password string
		expected func(ctrl *controller.Controller, actual int, err error)
	}{
		{
			name:     "login and logout test",
			userID:   1,
			login:    "test_user",
			password: "password1",
		},
	}

	for _, test := range tests {
		tt := test
		t.WithNewStep(tt.name, func(ctx provider.StepCtx) {
			ctrl := controller.NewController(t)

			authRepository := authRepo.New(ctrl.GetRedis())

			// create session
			sessionID := authRepository.CreateSession(tt.userID)
			ctx.Assert().Greater(len(sessionID), 0)

			// check session
			ctx.Assert().True(authRepository.SessionExists(sessionID))

			// get userID by sessionID
			userID, err := authRepository.GetUserIDBySessionID(sessionID)
			ctx.Assert().NoError(err)
			ctx.Assert().Equal(tt.userID, userID)

			// delete session
			err = authRepository.DeleteSession(sessionID)
			ctx.Assert().NoError(err)

			// check session
			ctx.Assert().False(authRepository.SessionExists(sessionID))
		})
	}
}
