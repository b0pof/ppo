package repository_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	userRepo "github.com/b0pof/ppo/internal/repository/user"
	"github.com/b0pof/ppo/tests/controller"
	"github.com/b0pof/ppo/tests/integration/common/creator/user"
)

type RepoUserFlow struct {
	suite.Suite
}

func (g *RepoUserFlow) TestUser(t provider.T) {
	t.Run("user flow", func(t provider.T) {
		ctrl := controller.NewController(t)
		// defer ctrl.Finish()

		userRepository := userRepo.New(ctrl.GetDB())

		ctx := context.Background()

		testUser := user.New(ctrl).Random()

		// get user by id
		user, err := userRepository.GetByID(ctx, testUser.ID)
		t.Assert().NoError(err)
		t.Assert().Equal(model.User{
			ID:       testUser.ID,
			Name:     testUser.Name,
			Login:    testUser.Login,
			Role:     testUser.Role,
			Phone:    testUser.Phone,
			Password: testUser.Password,
		}, user)

		// get user by login
		user, err = userRepository.GetByLogin(ctx, testUser.Login)
		t.Assert().NoError(err)
		t.Assert().Equal(model.User{
			ID:       testUser.ID,
			Name:     testUser.Name,
			Login:    testUser.Login,
			Role:     testUser.Role,
			Phone:    testUser.Phone,
			Password: testUser.Password,
		}, user)

		// get role
		role, err := userRepository.GetRoleByID(ctx, testUser.ID)
		t.Assert().NoError(err)
		t.Assert().Equal(model.RoleBuyer, role)

		newLogin := uuid.New().String()

		// update user
		err = userRepository.UpdateByID(ctx, testUser.ID, model.User{
			Login: newLogin,
			Name:  newLogin,
		})
		t.Assert().NoError(err)

		// check updated user
		user, err = userRepository.GetByID(ctx, testUser.ID)
		t.Assert().NoError(err)

		t.Assert().Equal(user.Login, newLogin)
		t.Assert().Equal(user.Name, newLogin)
	})
}
