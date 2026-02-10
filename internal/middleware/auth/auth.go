package auth

import (
	"errors"
	"net/http"

	"github.com/b0pof/ppo/internal/config"
	"github.com/b0pof/ppo/internal/util/pointer"
	"github.com/gorilla/mux"

	authUtil "github.com/b0pof/ppo/internal/util/auth"
	"github.com/b0pof/ppo/internal/util/cookie"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func New(auth auth, ver verification, user user, mode string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && mode == config.ServerModeReadOnly {
				response.Forbidden(w, pointer.To("this replica supports GET method only"))
				return
			}

			ctx := r.Context()

			// role := model.RoleGuest

			sessionID, err := cookie.GetSession(r)
			if !errors.Is(err, http.ErrNoCookie) && err != nil {
				response.BadRequest(w, err.Error())
				return
			}

			userID, _ := auth.GetUserIDBySessionID(sessionID)
			ctx = authUtil.WithUserID(ctx, userID)

			role, _ := user.GetRoleByID(ctx, userID)
			ctx = authUtil.WithRole(ctx, role)

			status, _ := ver.GetStatus(ctx, userID)
			ctx = authUtil.WithVerificationStatus(ctx, status)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
