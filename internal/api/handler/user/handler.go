package user

import usecase "github.com/b0pof/ppo/internal/usecase/user"

type User struct {
	user usecase.IUserUsecase
	log  logger
}

func New(u usecase.IUserUsecase, l logger) *User {
	return &User{
		user: u,
		log:  l,
	}
}
