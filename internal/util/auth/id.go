package auth

import (
	"context"
)

type userIdKey struct{}

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIdKey{}, userID)
}

func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(userIdKey{}).(int64)
	if !ok {
		return 0
	}

	return userID
}
