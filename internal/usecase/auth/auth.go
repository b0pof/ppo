package auth

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	passwordUtil "git.iu7.bmstu.ru/kia22u475/ppo/internal/util/password"
)

type IAuthUsecase interface {
	Login(ctx context.Context, login string, password string) (string, error)
	Signup(ctx context.Context, login string, password string, role string) (string, error)
	Logout(sessionID string) error
	IsLoggedIn(sessionID string) bool
	GetUserIDBySessionID(sessionID string) (int64, error)
}

type Usecase struct {
	auth authRepo
	user userRepo
}

func New(a authRepo, u userRepo) *Usecase {
	return &Usecase{
		auth: a,
		user: u,
	}
}

func (u *Usecase) Login(ctx context.Context, login string, password string) (string, error) {
	if login == "" || password == "" {
		return "", model.ErrInvalidInput
	}

	user, err := u.user.GetByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("auth usecase: failed to get user password: %w", err)
	}

	if !passwordUtil.Equal(password, user.Password) {
		return "", model.ErrWrongPassword
	}

	return u.auth.CreateSession(user.ID), nil
}

func (u *Usecase) Signup(ctx context.Context, login string, password string, role string) (string, error) {
	if login == "" || password == "" || role == "" {
		return "", model.ErrInvalidInput
	}

	passwordHash, err := passwordUtil.Hash(password)
	if err != nil {
		return "", model.ErrFailedToHash
	}

	userID, err := u.user.Create(ctx, login, passwordHash, role)
	if err != nil {
		return "", fmt.Errorf("auth usecase: failed to create user: %w", err)
	}

	sessionID := u.auth.CreateSession(userID)

	return sessionID, nil
}

func (u *Usecase) Logout(sessionID string) error {
	err := u.auth.DeleteSession(sessionID)
	if err != nil {
		return fmt.Errorf("auth usecase: failed to delete session: %w", err)
	}

	return nil
}

func (u *Usecase) IsLoggedIn(sessionID string) bool {
	return u.auth.SessionExists(sessionID)
}

func (u *Usecase) GetUserIDBySessionID(sessionID string) (int64, error) {
	userID, err := u.auth.GetUserIDBySessionID(sessionID)
	if err != nil {
		return 0, fmt.Errorf("auth usecase: failed to get user id: %w", err)
	}

	return userID, nil
}
