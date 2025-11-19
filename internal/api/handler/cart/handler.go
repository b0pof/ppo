package cart

import usecase "github.com/b0pof/ppo/internal/usecase/cart"

type Cart struct {
	cart usecase.ICartUsecase
	log  logger
}

func New(c usecase.ICartUsecase, l logger) *Cart {
	return &Cart{
		cart: c,
		log:  l,
	}
}
