package user

import (
	"context"
	"fmt"
	"time"

	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/repository/verification"
	passwordUtil "github.com/b0pof/ppo/internal/util/password"
)

type IUserUsecase interface {
	GetByID(ctx context.Context, userID int64) (model.User, error)
	GetRoleByID(ctx context.Context, userID int64) (string, error)
	GetMetaInfoByUserID(ctx context.Context, userID int64) (model.UserMetaInfo, error)
	UpdateByID(ctx context.Context, userID int64, user model.User) error
	UpdatePassword(ctx context.Context, userID int64, oldPassword string, password string, verified bool) error
}

type Usecase struct {
	user userRepo
	auth authUsecase
}

func New(u userRepo, a authUsecase) *Usecase {
	return &Usecase{
		user: u,
		auth: a,
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

func (u *Usecase) UpdatePassword(ctx context.Context, userID int64, oldPassword string, password string, verified bool) error {
	err := model.ValidateUserPassword(password)
	if err != nil {
		return fmt.Errorf("user usecase: failed to validate password: %w", err)
	}

	user, err := u.user.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user usecase: failed to get user: %w", err)
	}

	state, err := u.getState(ctx, oldPassword, user)
	if err != nil {
		return fmt.Errorf("user usecase: failed to get user state: %w", err)
	}

	if !state.ExpiresAt.Add(6*time.Hour).After(time.Now()) && state.Success {
		err = u.user.SaveTempPasswords(ctx, userID, oldPassword, password)
		if err != nil {
			return fmt.Errorf("user usecase: failed to save temp passwords: %w", err)
		}

		_, err = u.auth.SendCode(ctx, userID, user.Login, model.VerificationPurposePasswordChange)
		if err != nil {
			return fmt.Errorf("user usecase: failed to send code: %w", err)
		}
		return model.ErrNeedToVerify
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

func (u *Usecase) getState(ctx context.Context, oldPassword string, user model.User) (verification.VerificationCodeRow, error) {
	if !passwordUtil.Equal(oldPassword, user.Password) {
		return verification.VerificationCodeRow{}, model.ErrWrongPassword
	}

	state, err := verification.GlobalVerificationRepo.FetchAttempts(ctx, user.ID)
	if err != nil {
		return verification.VerificationCodeRow{}, fmt.Errorf("user usecase: failed to fetch verification state: %w", err)
	}

	return state, nil
}
