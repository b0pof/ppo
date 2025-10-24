package user

import usecase "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/user"

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
