package auth

import (
	"context"

	"github.com/b0pof/ppo/internal/model"
)

type verificationStatusKey struct{}

func WithVerificationStatus(ctx context.Context, state string) context.Context {
	return context.WithValue(ctx, verificationStatusKey{}, state)
}

func GetVerificationState(ctx context.Context) string {
	state, ok := ctx.Value(verificationStatusKey{}).(string)
	if !ok || state == "" {
		return model.VerificationStateFailed
	}

	return state
}
