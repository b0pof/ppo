package auth

import (
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
)

func (r *Repository) DeleteSession(sID string) error {
	err := r.redis.Get(sID).Err()
	if err != nil {
		return model.ErrNotFound
	}

	r.redis.Del(sID)

	return nil
}
