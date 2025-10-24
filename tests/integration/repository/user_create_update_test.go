//go:build integration

package repository

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	userRepo "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/user"
	"git.iu7.bmstu.ru/kia22u475/ppo/tests/controller"
)

type UserFlow struct {
	suite.Suite
}

func (g *UserFlow) TestUser(t provider.T) {
	tests := []struct {
		name   string
		userID int64
		login  string
	}{
		{
			name:   "order test",
			userID: 1,
			login:  "user1",
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t provider.T) {
			ctrl := controller.NewController(t)

			userRepository := userRepo.New(ctrl.GetDB())

			ctx := context.Background()

			// get user by id
			user, err := userRepository.GetByID(ctx, tt.userID)
			t.Assert().NoError(err)
			t.Assert().Equal(model.User{
				ID:       tt.userID,
				Name:     "User1 Name",
				Login:    tt.login,
				Role:     model.RoleBuyer,
				Phone:    "88005553535",
				Password: "$2a$10$fMfBQXaBUDqfl5VU2xOmpeUtX.l.Vx1ZfkaDADAcnzr7hHYZqPru.",
			}, user)

			// get user by login
			user, err = userRepository.GetByLogin(ctx, tt.login)
			t.Assert().NoError(err)
			t.Assert().Equal(model.User{
				ID:       tt.userID,
				Name:     "User1 Name",
				Login:    tt.login,
				Role:     model.RoleBuyer,
				Phone:    "88005553535",
				Password: "$2a$10$fMfBQXaBUDqfl5VU2xOmpeUtX.l.Vx1ZfkaDADAcnzr7hHYZqPru.",
			}, user)

			// get role
			role, err := userRepository.GetRoleByID(ctx, tt.userID)
			t.Assert().NoError(err)
			t.Assert().Equal(model.RoleBuyer, role)

			// create user
			newUserID, err := userRepository.Create(ctx, "newUser", "password123", model.RoleBuyer)
			t.Assert().NoError(err)

			expectedUser := model.User{
				ID:       newUserID,
				Name:     "newUser",
				Login:    "newUser",
				Role:     model.RoleBuyer,
				Password: "password123",
			}

			// get created user
			user, err = userRepository.GetByID(ctx, newUserID)
			t.Assert().NoError(err)
			t.Assert().Equal(expectedUser, user)

			// update user
			err = userRepository.UpdateByID(ctx, newUserID, model.User{
				Login: "updatedLogin",
				Name:  "updatedName",
			})
			t.Assert().NoError(err)

			// check updated user
			user, err = userRepository.GetByID(ctx, newUserID)
			t.Assert().NoError(err)

			expectedUser.Login = "updatedLogin"
			expectedUser.Name = "updatedName"
			t.Assert().Equal(expectedUser, user)
		})

	}
}
