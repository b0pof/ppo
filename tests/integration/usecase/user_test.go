package usecase_test

//
//import (
//	"context"
//
//	"github.com/Pallinder/go-randomdata"
//	"github.com/b0pof/ppo/tests/integration/common/generator/field/text"
//	"github.com/ozontech/allure-go/pkg/framework/provider"
//	"github.com/ozontech/allure-go/pkg/framework/suite"
//
//	"github.com/b0pof/ppo/internal/model"
//	authRepo "github.com/b0pof/ppo/internal/repository/auth"
//	userRepo "github.com/b0pof/ppo/internal/repository/user"
//	authUc "github.com/b0pof/ppo/internal/usecase/auth"
//	userUc "github.com/b0pof/ppo/internal/usecase/user"
//	"github.com/b0pof/ppo/tests/controller"
//)
//
//type UserTest struct {
//	suite.Suite
//}
//
//func (g *UserTest) TestFullUserFlow(t provider.T) {
//	var (
//		testPassword    = "pwdpwdpwd777"
//		newTestPassword = "newpwdpwd777"
//	)
//
//	t.WithNewStep("user flow", func(ctxA provider.StepCtx) {
//		ctx := context.Background()
//		ctrl := controller.NewController(t)
//		// defer ctrl.Finish()
//
//		userRepository := userRepo.New(ctrl.GetDB())
//		authRepository := authRepo.New(ctrl.GetRedis())
//
//		userUsecase := userUc.New(userRepository)
//		authUsecase := authUc.New(authRepository, userRepository)
//
//		testLogin := randomdata.SillyName()
//		_, err := authUsecase.Signup(ctx, testLogin, testPassword, model.RoleBuyer)
//		ctxA.Assert().NoError(err)
//
//		newUser, err := userRepository.GetByLogin(ctx, testLogin)
//		ctxA.Assert().NoError(err)
//
//		user, err := userUsecase.GetByID(ctx, newUser.ID)
//		ctxA.Assert().NoError(err)
//		ctxA.Assert().Equal(newUser.ID, user.ID)
//		ctxA.Assert().NotEmpty(user.Role)
//		ctxA.Assert().NotEmpty(user.Name)
//		ctxA.Assert().NotEmpty(user.Login)
//		ctxA.Assert().NotEmpty(user.Password)
//
//		role, err := userUsecase.GetRoleByID(ctx, newUser.ID)
//		ctxA.Assert().NoError(err)
//		ctxA.Assert().True(role == "buyer" || role == "seller" || role == "admin")
//
//		meta, err := userUsecase.GetMetaInfoByUserID(ctx, newUser.ID)
//		ctxA.Assert().NoError(err)
//		ctxA.Assert().Equal(user.Name, meta.Name)
//
//		updatedUser := model.User{
//			ID:    user.ID,
//			Name:  "Mark",
//			Login: text.RandStringRunes(10),
//		}
//		err = userUsecase.UpdateByID(ctx, newUser.ID, updatedUser)
//		ctxA.Assert().NoError(err)
//
//		afterUpdate, err := userUsecase.GetByID(ctx, newUser.ID)
//		ctxA.Assert().NoError(err)
//		ctxA.Assert().Equal(updatedUser.Name, afterUpdate.Name)
//		ctxA.Assert().Equal(updatedUser.Login, afterUpdate.Login)
//
//		err = userUsecase.UpdatePassword(ctx, newUser.ID, "wrongPass", newTestPassword)
//		ctxA.Assert().Error(err)
//		ctxA.Assert().ErrorIs(err, model.ErrWrongPassword)
//
//		err = userUsecase.UpdatePassword(ctx, newUser.ID, testPassword, newTestPassword)
//		ctxA.Assert().NoError(err)
//
//		userAfterPass, err := userUsecase.GetByID(ctx, newUser.ID)
//		ctxA.Assert().NoError(err)
//		ctxA.Assert().NotEqual(user.Password, userAfterPass.Password)
//	})
//}
