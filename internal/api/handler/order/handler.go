package order

import usecase "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/order"

type Order struct {
	order usecase.IOrderUsecase
	log   logger
}

func New(o usecase.IOrderUsecase, l logger) *Order {
	return &Order{
		order: o,
		log:   l,
	}
}
