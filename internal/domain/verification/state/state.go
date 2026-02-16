package state

import (
	"time"

	verificationCodeAttempts "github.com/b0pof/ppo/internal/domain/verification/code/attempts"
)

type VerificationState struct {
	Code         string
	Attempts     int64
	ExpiresAt    time.Time
	BlockedUntil *time.Time
	Success      bool
}

func New(code string, attempts int64, expiresAt time.Time, blockedUntil *time.Time, success bool) *VerificationState {
	return &VerificationState{
		Code:         code,
		Attempts:     attempts,
		ExpiresAt:    expiresAt,
		BlockedUntil: blockedUntil,
		Success:      success,
	}
}

func (s *VerificationState) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *VerificationState) IsBlocked() bool {
	if s.BlockedUntil == nil {
		return false
	}

	return s.BlockedUntil.After(time.Now())
}

func (s *VerificationState) IsPassed() bool {
	return s.Success
}

func (s *VerificationState) IsAvailable() bool {
	return !s.IsExpired() && !s.IsBlocked()
}

func (s *VerificationState) IsLimitReached() bool {
	return verificationCodeAttempts.MaxReached(s.Attempts)
}
