//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package auth

import (
	"context"

	"github.com/b0pof/ppo/internal/model"
)

type authRepo interface {
	CreateSession(userID int64) string
	SessionExists(sessionID string) bool
	GetUserIDBySessionID(sessionID string) (int64, error)
	DeleteSession(sessionID string) error
}

type userRepo interface {
	GetByLogin(ctx context.Context, login string) (model.User, error)
	Create(ctx context.Context, login string, password string, role string) (int64, error)
}
