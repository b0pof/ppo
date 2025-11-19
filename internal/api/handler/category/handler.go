package category

import (
	"github.com/b0pof/ppo/internal/repository/category"
	"github.com/b0pof/ppo/internal/repository/item"
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
