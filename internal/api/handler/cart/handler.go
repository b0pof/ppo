package cart

import usecase "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/cart"

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
