package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/b0pof/ppo/internal/model"
)

func (r *Repository) GetByID(ctx context.Context, userID int64) (model.User, error) {
	q := `
		select id, name, login, role, phone, password
		from "user"
		where id = $1;
	`

	var user userRow

	err := r.db.GetContext(ctx, &user, q, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model.ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("user repository.GetByID: %w", err)
	}

	return convertUser(user), nil
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	q := `
		select id, name, login, role, phone, password
		from "user"
		where login = $1;
	`

	var user userRow

	err := r.db.GetContext(ctx, &user, q, login)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model.ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("user repository.GetByLogin: %w", err)
	}

	return convertUser(user), nil
}

func (r *Repository) GetUserMetaByID(ctx context.Context, userID int64) (model.UserMetaInfo, error) {
	q := `
		select u.name, sum(case when ci.item_id is not null then 1 else 0 end) as cart_items_amount
		from "user" u
		left join cart_item ci on u.id = ci.user_id
		where u.id = $1
		group by u.name;
	`

	var userMeta userMetaRow

	err := r.db.GetContext(ctx, &userMeta, q, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.UserMetaInfo{}, model.ErrNotFound
	}
	if err != nil {
		return model.UserMetaInfo{}, fmt.Errorf("user repository.GetUserMetaByID: %w", err)
	}

	return convertUserMeta(userMeta), nil
}

func (r *Repository) GetUserLoginByID(ctx context.Context, userID int64) (string, error) {
	q := `
		select login
		from "user"
		where id = $1;
	`

	var login string

	err := r.db.GetContext(ctx, &login, q, userID)
	if err != nil {
		return "", fmt.Errorf("user repository.GetUserLoginByID: %w", err)
	}

	return login, nil
}

func (r *Repository) GetRoleByID(ctx context.Context, userID int64) (string, error) {
	q := `
		select role
		from "user"
		where id = $1;
	`

	var role string

	err := r.db.GetContext(ctx, &role, q, userID)
	if err != nil {
		return "", fmt.Errorf("user repository.GetRoleByID: %w", err)
	}

	return role, nil
}
