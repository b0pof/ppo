//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package user

import (
	"context"

	"github.com/b0pof/ppo/internal/model"
)

type userRepo interface {
	GetByID(ctx context.Context, userID int64) (model.User, error)
	GetUserMetaByID(ctx context.Context, userID int64) (model.UserMetaInfo, error)
	GetUserLoginByID(ctx context.Context, userID int64) (string, error)
	UpdatePasswordByID(ctx context.Context, userID int64, newPasswordHash string) error
	UpdateByID(ctx context.Context, userID int64, user model.User) error
	GetRoleByID(ctx context.Context, userID int64) (string, error)
}
