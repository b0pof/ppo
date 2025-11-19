package user

import (
	"context"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) Create(ctx context.Context, login string, password string, role string) (int64, error) {
	q := `
		WITH inserted AS (
			INSERT INTO "user" (name, login, password, role)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (login) DO NOTHING
			RETURNING id
		)
		SELECT COALESCE((SELECT id FROM inserted), 0) AS id;
	`

	var id int64

	err := r.db.GetContext(ctx, &id, q, login, login, password, role)
	if err != nil {
		return 0, fmt.Errorf("user repository.Create: %w", err)
	}

	if id == 0 {
		return 0, model.ErrAlreadyExists
	}

	return id, nil
}
