package auth

import (
	"context"

	"github.com/b0pof/ppo/internal/model"
)

type roleKey struct{}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey{}, role)
}

func GetRole(ctx context.Context) string {
	role, ok := ctx.Value(roleKey{}).(string)
	if !ok || role == "" {
		return model.RoleGuest
	}

	return role
}
