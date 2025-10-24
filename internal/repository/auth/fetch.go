package auth

import "git.iu7.bmstu.ru/kia22u475/ppo/internal/model"

func (r *Repository) GetUserIDBySessionID(sID string) (int64, error) {
	uID, err := r.redis.Get(sID).Int64()
	if err != nil {
		return 0, model.ErrNotFound
	}

	return uID, nil
}

func (r *Repository) SessionExists(sID string) bool {
	uID, err := r.redis.Get(sID).Int64()
	if err != nil {
		return false
	}

	return uID != 0
}
