//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package auth

import "context"

type auth interface {
	GetUserIDBySessionID(sessionID string) (int64, error)
}

type user interface {
	GetRoleByID(ctx context.Context, userID int64) (string, error)
}
