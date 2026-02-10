//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package auth

import (
	"context"
	"time"

	"github.com/b0pof/ppo/internal/model"
	verificationRepository "github.com/b0pof/ppo/internal/repository/verification"
)

type authRepo interface {
	CreateSession(userID int64) string
	SessionExists(sessionID string) bool
	GetUserIDBySessionID(sessionID string) (int64, error)
	DeleteSession(sessionID string) error
}

type verificationRepo interface {
	Create(ctx context.Context, req model.SaveVerificationCodeRequest) error
	FetchAttempts(ctx context.Context, userID int64) (verificationRepository.VerificationCodeRow, error)
	IncrementAttempts(ctx context.Context, userID int64, blockedUntil *time.Time) error
	SetStatus(ctx context.Context, userID int64, success bool) error
}

type userRepo interface {
	GetByLogin(ctx context.Context, login string) (model.User, error)
	Create(ctx context.Context, login string, password string, role string) (int64, error)
	GetTempPasswords(ctx context.Context, userID int64) (string, string, error)
}

type userUsecase interface {
	UpdatePassword(ctx context.Context, userID int64, oldPassword string, password string, verified bool) error
}

type notifierUsecase interface {
	SendCode(ctx context.Context, email string, code string) error
}
