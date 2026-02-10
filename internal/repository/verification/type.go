package verification

import "time"

type VerificationCodeRow struct {
	UserID       int64      `db:"user_id"`
	Code         string     `db:"code"`
	Attempts     int64      `db:"attempts"`
	ExpiresAt    time.Time  `db:"expires_at"`
	BlockedUntil *time.Time `db:"blocked_until"`
	Purpose      string     `db:"purpose"`
}
