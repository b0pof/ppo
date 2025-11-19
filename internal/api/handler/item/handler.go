package item

import usecase "github.com/b0pof/ppo/internal/usecase/item"

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
