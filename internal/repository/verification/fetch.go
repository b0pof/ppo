package verification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/b0pof/ppo/internal/config"
	"github.com/b0pof/ppo/internal/configure"
	"github.com/b0pof/ppo/internal/model"
)

var GlobalVerificationRepo = Repository{}

func init() {
	GlobalVerificationRepo = Repository{
		db: configure.MustInitPostgres(context.Background(), config.MustLoad().Postgres),
	}
}

func (r *Repository) FetchAttempts(ctx context.Context, userID int64) (VerificationCodeRow, error) {
	q := `
		select user_id, code, attempts, expires_at, blocked_until, purpose, success
		from verification_code
		where user_id = $1;
	`

	var verification VerificationCodeRow

	err := r.db.GetContext(ctx, &verification, q, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return VerificationCodeRow{}, model.ErrNotFound
	}
	if err != nil {
		return VerificationCodeRow{}, fmt.Errorf("verification repository.FetchAttempts: %w", err)
	}

	return verification, nil
}

func (r *Repository) GetStatus(ctx context.Context, userID int64) (string, error) {
	q := `
		select success
		from verification_code
		where user_id = $1
			and updated_at > now() - interval '30 min';
	`

	var success bool

	err := r.db.GetContext(ctx, &success, q, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.VerificationStateFailed, nil
	}
	if err != nil {
		return "", fmt.Errorf("verification repository.GetStatus: %w", err)
	}

	status := model.VerificationStateFailed
	if success {
		status = model.VerificationStateVerified
	}

	fmt.Println("status:", status)

	return status, nil
}
