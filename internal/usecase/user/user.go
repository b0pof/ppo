package user

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	passwordUtil "git.iu7.bmstu.ru/kia22u475/ppo/internal/util/password"
)

type IUserUsecase interface {
	GetByID(ctx context.Context, userID int64) (model.User, error)
	GetRoleByID(ctx context.Context, userID int64) (string, error)
	GetMetaInfoByUserID(ctx context.Context, userID int64) (model.UserMetaInfo, error)
	UpdateByID(ctx context.Context, userID int64, user model.User) error
	UpdatePassword(ctx context.Context, userID int64, oldPassword string, password string) error
}

type Usecase struct {
	user userRepo
}

func New(u userRepo) *Usecase {
	return &Usecase{
		user: u,
	}
}

func (u *Usecase) GetByID(ctx context.Context, userID int64) (model.User, error) {
	user, err := u.user.GetByID(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("user usecase: failed to get user: %w", err)
	}

	return user, nil
}

func (u *Usecase) GetRoleByID(ctx context.Context, userID int64) (string, error) {
	role, err := u.user.GetRoleByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user usecase: failed to get role: %w", err)
	}

	return role, nil
}

func (u *Usecase) GetMetaInfoByUserID(ctx context.Context, userID int64) (model.UserMetaInfo, error) {
	userMeta, err := u.user.GetUserMetaByID(ctx, userID)
	if err != nil {
		return model.UserMetaInfo{}, fmt.Errorf("user usecase: failed to get user meta: %w", err)
	}

	return userMeta, nil
}

func (u *Usecase) UpdateByID(ctx context.Context, userID int64, user model.User) error {
	err := model.ValidateUser(user)
	if err != nil {
		return fmt.Errorf("user usecase: failed to validate user: %w", err)
	}

	err = u.user.UpdateByID(ctx, userID, user)
	if err != nil {
		return fmt.Errorf("user usecase: failed to update user: %w", err)
	}

	return nil
}

func (u *Usecase) UpdatePassword(ctx context.Context, userID int64, oldPassword string, password string) error {
	err := model.ValidateUserPassword(password)
	if err != nil {
		return fmt.Errorf("user usecase: failed to validate password: %w", err)
	}

	user, err := u.user.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user usecase: failed to get user: %w", err)
	}

	if !passwordUtil.Equal(oldPassword, user.Password) {
		return model.ErrWrongPassword
	}

	newPassword, err := passwordUtil.Hash(password)
	if err != nil {
		return model.ErrFailedToHash
	}

	err = u.user.UpdatePasswordByID(ctx, userID, newPassword)
	if err != nil {
		return fmt.Errorf("user usecase: failed to update password: %w", err)
	}

	return nil
}
