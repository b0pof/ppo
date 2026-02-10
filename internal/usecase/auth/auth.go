package auth

import (
	"context"
	"fmt"
	"time"

	verificationBlock "github.com/b0pof/ppo/internal/domain/verification/block"
	verificationCodeAttempts "github.com/b0pof/ppo/internal/domain/verification/code/attempts"
	verificationCodeExpiration "github.com/b0pof/ppo/internal/domain/verification/code/expiration"
	verificationCodeGenerator "github.com/b0pof/ppo/internal/domain/verification/code/generator"
	verificationState "github.com/b0pof/ppo/internal/domain/verification/state"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/repository/verification"
	authUtil "github.com/b0pof/ppo/internal/util/auth"
	passwordUtil "github.com/b0pof/ppo/internal/util/password"
	"github.com/b0pof/ppo/internal/util/pointer"
)

type IAuthUsecase interface {
	Login(ctx context.Context, login string, password string) (*string, error)
	Signup(ctx context.Context, login string, password string, role string) error
	Logout(sessionID string) error
	IsLoggedIn(sessionID string) bool
	GetUserIDBySessionID(sessionID string) (int64, error)
	ApplyVerificationCode(ctx context.Context, email string, code string) (model.ApplyVerificationCodeResult, error)
}

type Usecase struct {
	auth         authRepo
	verification verificationRepo
	user         userRepo
	userManager  userUsecase
	notification notifierUsecase
}

func (u *Usecase) SetUserUsecase(uu userUsecase) {
	u.userManager = uu
}

func New(a authRepo, v verificationRepo, u userRepo, n notifierUsecase) *Usecase {
	return &Usecase{
		auth:         a,
		verification: v,
		user:         u,
		notification: n,
	}
}

func (u *Usecase) Login(ctx context.Context, login string, password string) (*string, error) {
	if login == "" || password == "" {
		return nil, model.ErrInvalidInput
	}

	user, err := u.user.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("auth usecase: failed to get user password: %w", err)
	}

	if !passwordUtil.Equal(password, user.Password) {
		return nil, model.ErrWrongPassword
	}

	if authUtil.GetVerificationState(ctx) == model.VerificationStateVerified {
		return pointer.To(u.auth.CreateSession(user.ID)), nil
	}

	status, err := u.verification.FetchAttempts(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("auth usecase: Login: failed to get verification status: %w", err)
	}

	state := verificationState.New(status.Code, status.Attempts, status.ExpiresAt, status.BlockedUntil)

	if state.IsBlocked() {
		return nil, model.ErrTemporaryBlocked
	}

	_, err = u.SendCode(ctx, user.ID, user.Login, model.VerificationPurposeAuth)
	if err != nil {
		return nil, fmt.Errorf("auth usecase: failed to send code: %w", err)
	}

	return nil, nil
}

func (u *Usecase) ApplyVerificationCode(ctx context.Context, email string, code string) (model.ApplyVerificationCodeResult, error) {
	user, status, err := u.getStatusAndUser(ctx, email, code)
	if err != nil {
		return model.ApplyVerificationCodeResult{}, err
	}

	state := verificationState.New(status.Code, status.Attempts, status.ExpiresAt, status.BlockedUntil)

	if state.IsBlocked() || state.IsLimitReached() {
		return model.ApplyVerificationCodeResult{
			Success:        false,
			Message:        "Попытки закончились",
			RetryAvailable: pointer.To(false),
			RetryAfter:     pointer.To(int(time.Until(*status.BlockedUntil).Seconds())),
		}, nil
	}

	if state.IsExpired() {
		expiresAt, err := u.SendCode(ctx, user.ID, user.Login, status.Purpose)
		return model.ApplyVerificationCodeResult{
			Success:        false,
			Message:        "Новый код подтверждения отправлен на почту",
			ExpiresIn:      pointer.To(int(time.Until(expiresAt).Seconds())),
			RetryAvailable: pointer.To(true),
		}, err
	}

	isCodeCorrect := code == status.Code
	if err = u.verification.SetStatus(ctx, user.ID, isCodeCorrect); err != nil {
		return model.ApplyVerificationCodeResult{}, fmt.Errorf("auth usecase: failed to set verification status: %w", err)
	}

	if !isCodeCorrect {
		message := "Неверный код"
		var blockUntil *time.Time

		if verificationCodeAttempts.MaxReached(status.Attempts + 1) {
			blockUntil = pointer.To(verificationBlock.New().GetExpiration())
			message = "Попытки закончились"
		}

		err = u.verification.IncrementAttempts(ctx, user.ID, blockUntil)
		return model.ApplyVerificationCodeResult{
			Success:        false,
			Message:        message,
			ExpiresIn:      pointer.To(int(time.Until(status.ExpiresAt).Seconds())),
			RetryAvailable: pointer.To(blockUntil == nil),
		}, err
	}

	if status.Purpose == model.VerificationPurposePasswordChange {
		if err := u.updatePassword(ctx, user.ID); err != nil {
			return model.ApplyVerificationCodeResult{}, err
		}
	}

	return model.ApplyVerificationCodeResult{
		Success:   true,
		Message:   "Успех",
		SessionID: pointer.To(u.auth.CreateSession(user.ID)),
	}, nil
}

func (u *Usecase) getStatusAndUser(ctx context.Context, email, code string) (model.User, verification.VerificationCodeRow, error) {
	if email == "" || code == "" {
		return model.User{}, verification.VerificationCodeRow{}, model.ErrInvalidInput
	}

	user, err := u.user.GetByLogin(ctx, email)
	if err != nil {
		return model.User{}, verification.VerificationCodeRow{}, fmt.Errorf("auth usecase: failed to get user password: %w", err)
	}

	status, err := u.verification.FetchAttempts(ctx, user.ID)
	if err != nil {
		return model.User{}, verification.VerificationCodeRow{}, fmt.Errorf("auth usecase: failed to get verification status: %w", err)
	}

	return user, status, nil
}

func (u *Usecase) updatePassword(ctx context.Context, userID int64) error {
	oldPwd, newPwd, err := u.user.GetTempPasswords(ctx, userID)
	if err != nil {
		return fmt.Errorf("auth usecase: failed to get temp passwords: %w", err)
	}

	err = u.userManager.UpdatePassword(ctx, userID, oldPwd, newPwd, true)
	if err != nil {
		return fmt.Errorf("auth usecase: failed to update password: %w", err)
	}

	return nil
}

func (u *Usecase) Signup(ctx context.Context, login string, password string, role string) error {
	if login == "" || password == "" || role == "" {
		return model.ErrInvalidInput
	}

	passwordHash, err := passwordUtil.Hash(password)
	if err != nil {
		return model.ErrFailedToHash
	}

	userID, err := u.user.Create(ctx, login, passwordHash, role)
	if err != nil {
		return fmt.Errorf("auth usecase: failed to create user: %w", err)
	}

	_, err = u.SendCode(ctx, userID, login, model.VerificationPurposeAuth)
	if err != nil {
		return fmt.Errorf("auth usecase: failed to send code: %w", err)
	}

	return nil
}

func (u *Usecase) SendCode(ctx context.Context, userID int64, email string, purpose string) (time.Time, error) {
	code := verificationCodeGenerator.New().Generate()
	expiresAt := verificationCodeExpiration.New().Get()

	err := u.verification.Create(ctx, model.SaveVerificationCodeRequest{
		UserID:    userID,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("auth usecase: Login: failed to create verification code: %w", err)
	}

	err = u.notification.SendCode(ctx, email, code)
	if err != nil {
		return time.Time{}, fmt.Errorf("auth usecase: failed to send verification code: %w", err)
	}

	return expiresAt, nil
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
