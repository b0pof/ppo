package order

import usecase "github.com/b0pof/ppo/internal/usecase/order"

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
