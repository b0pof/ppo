//go:build integration

package usecase

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/internal/model"
	userRepo "github.com/b0pof/ppo/internal/repository/user"
	userUsecase "github.com/b0pof/ppo/internal/usecase/user"
	"github.com/b0pof/ppo/tests/controller"
)

type UserTest struct {
	suite.Suite
}

func (g *UserTest) TestFullUserFlow(t provider.T) {
	tests := []struct {
		name    string
		userID  int64
		oldPass string
		newPass string
	}{
		{
			name:    "full user usecase flow",
			userID:  1,
			oldPass: "testtest",
			newPass: "newStrongPass1",
		},
	}

	for _, test := range tests {
		tt := test
		t.WithNewStep(tt.name, func(ctxA provider.StepCtx) {
			ctx := context.Background()
			ctrl := controller.NewController(t)

			repo := userRepo.New(ctrl.GetDB())
			usecase := userUsecase.New(repo)

			user, err := usecase.GetByID(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(tt.userID, user.ID)
			ctxA.Assert().NotEmpty(user.Role)
			ctxA.Assert().NotEmpty(user.Name)
			ctxA.Assert().NotEmpty(user.Login)
			ctxA.Assert().NotEmpty(user.Phone)
			ctxA.Assert().NotEmpty(user.Password)

			role, err := usecase.GetRoleByID(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().True(role == "buyer" || role == "seller" || role == "admin")

			meta, err := usecase.GetMetaInfoByUserID(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(user.Name, meta.Name)

			updatedUser := model.User{
				ID:    user.ID,
				Name:  "Mark",
				Login: "updated",
			}
			err = usecase.UpdateByID(ctx, tt.userID, updatedUser)
			ctxA.Assert().NoError(err)

			afterUpdate, err := usecase.GetByID(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().Equal(updatedUser.Name, afterUpdate.Name)
			ctxA.Assert().Equal(updatedUser.Login, afterUpdate.Login)

			err = usecase.UpdatePassword(ctx, tt.userID, "wrongPass", tt.newPass)
			ctxA.Assert().Error(err)
			ctxA.Assert().ErrorIs(err, model.ErrWrongPassword)

			err = usecase.UpdatePassword(ctx, tt.userID, tt.oldPass, tt.newPass)
			ctxA.Assert().NoError(err)

			userAfterPass, err := usecase.GetByID(ctx, tt.userID)
			ctxA.Assert().NoError(err)
			ctxA.Assert().NotEqual(user.Password, userAfterPass.Password)
		})
	}
}
