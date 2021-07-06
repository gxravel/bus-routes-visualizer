package middleware

import (
	"context"
	"net/http"
	"strings"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizer"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizercontext"
)

// Auth searches user by token and adds his data to context.
func Auth(visualizer *visualizer.Visualizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token := getAuthToken(r)

			err := visualizer.VerifyToken(ctx, token)
			if err != nil {
				log.
					FromContext(ctx).
					WithStr("token", token).
					Debug("verify token")

				api.RespondError(ctx, w, err)
				return
			}

			ctx = context.WithValue(ctx, visualizercontext.TokenKey, token)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

const (
	// AuthHeader is a header used to find token of user.
	AuthHeader = "Authorization"
)

func getAuthToken(r *http.Request) string {
	tokens, ok := r.Header[AuthHeader]
	if ok {
		if len(tokens) > 0 {
			return strings.TrimPrefix(tokens[0], "Bearer ")
		}
	}

	return ""
}
