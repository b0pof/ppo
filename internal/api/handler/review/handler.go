package review

import "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/review"

type Review struct {
	review *review.Repository
	log    logger
}

func New(r *review.Repository, l logger) *Review {
	return &Review{
		review: r,
		log:    l,
	}
}
