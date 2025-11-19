package auth

import (
	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) DeleteSession(sID string) error {
	err := r.redis.Get(sID).Err()
	if err != nil {
		return model.ErrNotFound
	}

	r.redis.Del(sID)

	return nil
}
