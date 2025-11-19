package review

import "github.com/b0pof/ppo/internal/repository/review"

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
