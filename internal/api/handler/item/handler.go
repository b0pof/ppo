package item

import usecase "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/item"

type Item struct {
	item usecase.IItemUsecase
	log  logger
}

func New(i usecase.IItemUsecase, l logger) *Item {
	return &Item{
		item: i,
		log:  l,
	}
}
