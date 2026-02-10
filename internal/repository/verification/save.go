package verification

import (
	"context"
	"fmt"
	"time"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) Create(ctx context.Context, req model.SaveVerificationCodeRequest) error {
	q := `
		insert into verification_code (user_id, code, purpose, expires_at)
		values ($1, $2, $3, $4)
		on conflict (user_id)
		do update set
			code = excluded.code,
			purpose = excluded.purpose,
			expires_at = excluded.expires_at,
			attempts = 0,
			blocked_until = null,
			success = false,
			created_at = now(),
			updated_at = now();
	`

	_, err := r.db.ExecContext(ctx, q, req.UserID, req.Code, req.Purpose, req.ExpiresAt)
	if err != nil {
		return fmt.Errorf("verification repository.Create: %w", err)
	}

	return nil
}

func (r *Repository) IncrementAttempts(ctx context.Context, userID int64, blockedUntil *time.Time) error {
	q := `
		update verification_code
		set
			attempts = attempts + 1,
			blocked_until = $1,
			updated_at = now()
		where user_id = $2
	`

	_, err := r.db.ExecContext(ctx, q, blockedUntil, userID)
	if err != nil {
		return fmt.Errorf("verification repository.IncrementAttempts: %w", err)
	}

	return nil
}

func (r *Repository) SetStatus(ctx context.Context, userID int64, success bool) error {
	q := `
		update verification_code
		set
			success = $1,
			updated_at = now()
		where user_id = $2
	`

	_, err := r.db.ExecContext(ctx, q, success, userID)
	if err != nil {
		return fmt.Errorf("verification repository.SetStatus: %w", err)
	}

	return nil
}
