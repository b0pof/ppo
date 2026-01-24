package user

import (
	"context"

	"github.com/Pallinder/go-randomdata"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/b0pof/ppo/internal/model"
	authRepo "github.com/b0pof/ppo/internal/repository/auth"
	userRepo "github.com/b0pof/ppo/internal/repository/user"
	authUc "github.com/b0pof/ppo/internal/usecase/auth"
)

type controller interface {
	GetDB() *sqlx.DB
	GetRedis() *redis.Client
}

type User struct {
	user *userRepo.Repository
	auth *authRepo.Repository
}

func New(c controller) *User {
	userRepository := userRepo.New(c.GetDB())
	authRepository := authRepo.New(c.GetRedis())

	return &User{
		user: userRepository,
		auth: authRepository,
	}
}

func (u *User) Random() model.User {
	authUsecase := authUc.New(u.auth, u.user)

	testLogin := uuid.New().String()
	testPassword := randomdata.Digits(20)

	ctx := context.Background()

	_, _ = authUsecase.Signup(ctx, testLogin, testPassword, model.RoleBuyer)
	newUser, _ := u.user.GetByLogin(ctx, testLogin)

	return newUser
}
