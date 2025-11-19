package permission

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"

	"github.com/b0pof/ppo/internal/model"
	authUtil "github.com/b0pof/ppo/internal/util/auth"
	"github.com/b0pof/ppo/internal/util/http/response"
)

type Middleware struct {
	permissions map[model.Path]model.Permissions
}

func New() *Middleware {
	return &Middleware{
		permissions: make(map[model.Path]model.Permissions),
	}
}

func (m *Middleware) Register(path model.Path, perms model.Permissions) {
	m.permissions[path] = perms
}

func (m *Middleware) New() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			role := authUtil.GetRole(ctx)

			if !m.hasAccess(model.NewPath(r.URL.Path, r.Method), role) {
				response.Forbidden(w, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middleware) hasAccess(path model.Path, role string) bool {
	for pattern, perms := range m.permissions {
		if equalUrls(pattern.Url(), path.Url()) && pattern.Method() == path.Method() {
			return perms.HasAccess(role)
		}
	}

	return false
}

func equalUrls(pattern, actual string) bool {
	escaped := regexp.QuoteMeta(pattern)

	escaped = strings.ReplaceAll(escaped, `\{`, "{")
	escaped = strings.ReplaceAll(escaped, `\}`, "}")

	replaced := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(escaped, `[^/]+`)

	regexPattern := "^" + replaced + "$"

	r, err := regexp.Compile(regexPattern)
	if err != nil {
		return false
	}

	return r.MatchString(actual)
}
