package user

import (
	"context"

	"github.com/Pallinder/go-randomdata"
	"github.com/b0pof/ppo/internal/config"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/b0pof/ppo/internal/model"
	authRepo "github.com/b0pof/ppo/internal/repository/auth"
	userRepo "github.com/b0pof/ppo/internal/repository/user"
	verificationRepo "github.com/b0pof/ppo/internal/repository/verification"
	authUc "github.com/b0pof/ppo/internal/usecase/auth"
	notificationUc "github.com/b0pof/ppo/internal/usecase/notification/post/verification/code"
)

type controller interface {
	GetDB() *sqlx.DB
	GetRedis() *redis.Client
}

type User struct {
	user         *userRepo.Repository
	auth         *authRepo.Repository
	verification *verificationRepo.Repository
	notification *notificationUc.Usecase
}

func New(c controller) *User {
	userRepository := userRepo.New(c.GetDB())
	authRepository := authRepo.New(c.GetRedis())
	verificationRepository := verificationRepo.New(c.GetDB())
	notification := notificationUc.New(config.GlobalCfg.SMTP)

	return &User{
		user:         userRepository,
		auth:         authRepository,
		verification: verificationRepository,
		notification: notification,
	}
}

func (u *User) Random() model.User {
	authUsecase := authUc.New(u.auth, u.verification, u.user, u.notification)

	testLogin := uuid.New().String()
	testPassword := randomdata.Digits(20)

	ctx := context.Background()

	_ = authUsecase.Signup(ctx, testLogin, testPassword, model.RoleBuyer)
	newUser, _ := u.user.GetByLogin(ctx, testLogin)

	return newUser
}
