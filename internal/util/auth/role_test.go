package auth_test

import (
	"context"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"

	. "git.iu7.bmstu.ru/kia22u475/ppo/internal/util/auth"
)

func TestWithRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		role         string
		expectations func(t assert.TestingT, got string)
	}{
		{
			name: "success",
			role: model.RoleBuyer,
			expectations: func(t assert.TestingT, got string) {
				assert.Equal(t, model.RoleBuyer, got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := WithRole(context.Background(), tc.role)
			out := GetRole(ctx)

			tc.expectations(t, out)
		})
	}
}

func TestGetRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		role         string
		putRole      bool
		expectations func(t assert.TestingT, got string)
	}{
		{
			name:    "success",
			role:    model.RoleBuyer,
			putRole: true,
			expectations: func(t assert.TestingT, got string) {
				assert.Equal(t, model.RoleBuyer, got)
			},
		},
		{
			name:    "no role found",
			putRole: false,
			expectations: func(t assert.TestingT, got string) {
				assert.Equal(t, model.RoleGuest, got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tc.putRole {
				ctx = WithRole(ctx, tc.role)
			}

			out := GetRole(ctx)

			tc.expectations(t, out)
		})
	}
}
