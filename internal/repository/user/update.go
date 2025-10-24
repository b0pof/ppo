package user

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
)

const (
	pgStringField = "'%s'"
)

func (r *Repository) UpdatePasswordByID(ctx context.Context, userID int64, newPasswordHash string) error {
	q := `
		update "user"
		set password = $1
		where id = $2;
	`

	_, err := r.db.ExecContext(ctx, q, newPasswordHash, userID)
	if err != nil {
		return fmt.Errorf("user repository.UpdatePasswordByID: %w", err)
	}

	return nil
}

func (r *Repository) UpdateByID(ctx context.Context, userID int64, user model.User) error {
	q := `
		update "user" set
		    name = %s,
			login = %s,
			phone = %s
		where id = $1;
	`

	newName := "name"
	if user.Name != "" {
		newName = fmt.Sprintf(pgStringField, user.Name)
	}

	newLogin := "login"
	if user.Login != "" {
		newLogin = fmt.Sprintf(pgStringField, user.Login)
	}

	newPhone := "phone"
	if user.Phone != "" {
		newPhone = fmt.Sprintf(pgStringField, user.Phone)
	}

	q = fmt.Sprintf(q, newName, newLogin, newPhone)

	_, err := r.db.ExecContext(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("user repository.UpdateByID: %w", err)
	}

	return nil
}
