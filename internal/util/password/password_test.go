package password_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/b0pof/ppo/internal/util/password"
)

func TestHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		password     string
		expectations func(t assert.TestingT, got string, err error)
	}{
		{
			name:     "success",
			password: "testtest",
			expectations: func(t assert.TestingT, got string, err error) {
				assert.Greater(t, len(got), 0)
				assert.NoError(t, err)
			},
		},
		{
			name:     "failed to hash",
			password: strings.Repeat("testtest", 10),
			expectations: func(t assert.TestingT, got string, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := Hash(tc.password)

			tc.expectations(t, out, err)
		})
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		password     string
		targetHash   string
		expectations func(t assert.TestingT, got bool)
	}{
		{
			name:       "equal",
			password:   "testtest",
			targetHash: "$2a$10$szUxJPeVofgzLzJ7UAaT4e2EH2n7fuEye7HOrVUmyiFTOapN1e.vK",
			expectations: func(t assert.TestingT, got bool) {
				assert.True(t, got)
			},
		},
		{
			name:       "not equal",
			password:   "another_one",
			targetHash: "$2a$10$szUxJPeVofgzLzJ7UAaT4e2EH2n7fuEye7HOrVUmyiFTOapN1e.vK",
			expectations: func(t assert.TestingT, got bool) {
				assert.False(t, got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := Equal(tc.password, tc.targetHash)

			tc.expectations(t, out)
		})
	}
}
