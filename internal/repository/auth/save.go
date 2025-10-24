package auth

import (
	"github.com/google/uuid"
)

func (r *Repository) CreateSession(userID int64) string {
	sID := uuid.New().String()

	r.redis.Set(sID, userID, r.sessionTTL)

	return sID
}
