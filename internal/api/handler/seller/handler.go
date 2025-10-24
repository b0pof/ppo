package seller

import usecase "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/item"

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
