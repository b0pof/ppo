package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/b0pof/ppo/internal/util/auth"
)

func TestWithUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		expectations func(t assert.TestingT, got int64)
	}{
		{
			name:   "success",
			userID: 1,
			expectations: func(t assert.TestingT, got int64) {
				assert.Equal(t, int64(1), got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := WithUserID(context.Background(), tc.userID)
			out := GetUserID(ctx)

			tc.expectations(t, out)
		})
	}
}

func TestGetUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       int64
		putUserID    bool
		expectations func(t assert.TestingT, got int64)
	}{
		{
			name:      "success",
			userID:    1,
			putUserID: true,
			expectations: func(t assert.TestingT, got int64) {
				assert.Equal(t, int64(1), got)
			},
		},
		{
			name:      "userID not found",
			userID:    1,
			putUserID: false,
			expectations: func(t assert.TestingT, got int64) {
				assert.Equal(t, int64(0), got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tc.putUserID {
				ctx = WithUserID(context.Background(), tc.userID)
			}

			out := GetUserID(ctx)

			tc.expectations(t, out)
		})
	}
}
