package model

import (
	"errors"
	"time"
)

const (
	VerificationPurposePasswordChange = "password_change"
	VerificationPurposeAuth           = "auth"

	VerificationStateVerified = "verified"
	VerificationStateFailed   = "failed"
)

var ErrTemporaryBlocked = errors.New("blocked")

type ApplyVerificationCodeResult struct {
	Success        bool
	Message        string
	SessionID      *string
	ExpiresIn      *int
	RetryAvailable *bool
	RetryAfter     *int
}

type SaveVerificationCodeRequest struct {
	UserID    int64
	Code      string
	Purpose   string
	ExpiresAt time.Time
}
