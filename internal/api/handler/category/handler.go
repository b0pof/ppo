package category

import (
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/category"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/item"
)

type Category struct {
	category *category.Repository
	item     *item.Repository
	log      logger
}

func New(c *category.Repository, i *item.Repository, l logger) *Category {
	return &Category{
		category: c,
		item:     i,
		log:      l,
	}
}
