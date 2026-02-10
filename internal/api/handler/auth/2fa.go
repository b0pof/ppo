package auth

import (
	"errors"
	"net/http"
	"time"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/request"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Auth) PostApiAuthVerify2fa(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := request.ParseBody[dto.Verify2FARequest](r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	result, err := h.auth.ApplyVerificationCode(ctx, data.Email, data.Code)
	if errors.Is(err, model.ErrNotFound) {
		response.BadRequest(w, "Пользователь не найден")
		return
	}
	if err != nil {
		h.log.Warn("login failed", err)
		response.Internal(w, err)
		return
	}

	if result.SessionID != nil {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    *result.SessionID,
			HttpOnly: true,
			Expires:  time.Now().Add(h.sessionTTL),
			Path:     "/",
		}

		http.SetCookie(w, cookie)
	}

	response.OK(w, present(result))
}

func present(in model.ApplyVerificationCodeResult) dto.Verify2FAResponse {
	return dto.Verify2FAResponse{
		Success:        &in.Success,
		Message:        &in.Message,
		ExpiresIn:      in.ExpiresIn,
		RetryAfter:     in.RetryAfter,
		RetryAvailable: in.RetryAvailable,
		Session:        in.SessionID,
	}
}
