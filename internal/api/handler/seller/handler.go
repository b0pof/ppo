package seller

import usecase "github.com/b0pof/ppo/internal/usecase/item"

type Seller struct {
	item usecase.IItemUsecase
	log  logger
}

func New(i usecase.IItemUsecase, l logger) *Seller {
	return &Seller{
		item: i,
		log:  l,
	}
}
