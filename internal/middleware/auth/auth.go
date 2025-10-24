package auth

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	authUtil "git.iu7.bmstu.ru/kia22u475/ppo/internal/util/auth"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/cookie"
	"git.iu7.bmstu.ru/kia22u475/ppo/internal/util/http/response"
)

func New(auth auth, user user) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			role := model.RoleGuest

			sessionID, err := cookie.GetSession(r)
			if !errors.Is(err, http.ErrNoCookie) && err != nil {
				response.BadRequest(w, err.Error())
				return
			}

			userID, _ := auth.GetUserIDBySessionID(sessionID)
			ctx = authUtil.WithUserID(ctx, userID)

			role, _ = user.GetRoleByID(ctx, userID)
			ctx = authUtil.WithRole(ctx, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
