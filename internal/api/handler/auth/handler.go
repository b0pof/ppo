package auth

import (
	"time"

	usecase "github.com/b0pof/ppo/internal/usecase/auth"
)

type Auth struct {
	auth usecase.IAuthUsecase
	log  logger

	sessionTTL time.Duration
}

func New(u usecase.IAuthUsecase, l logger, sessionTTL time.Duration) *Auth {
	return &Auth{
		auth:       u,
		log:        l,
		sessionTTL: sessionTTL,
	}
}
