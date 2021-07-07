package middleware

import (
	"context"
	"net/http"
	"strings"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
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

			userID, err := visualizer.VerifyToken(ctx, token)
			if err != nil {
				log.
					FromContext(ctx).
					WithStr("token", token).
					Debug("verify token")

				api.RespondError(ctx, w, err)
				return
			}

			filter := dataprovider.
				NewPermissionFilter().
				ByUserIDs(userID).
				ByActions(r.Method + ":" + r.URL.Path[7:])

			if err := visualizer.CheckPermission(ctx, filter); err != nil {
				log.
					FromContext(ctx).
					WithErr(err).
					WithFields(
						"userID", userID,
						"actions", filter.Actions,
					).
					Debug("check permission")
				api.RespondError(ctx, w, err)
				return
			}

			ctx = context.WithValue(ctx, visualizercontext.TokenKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
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
