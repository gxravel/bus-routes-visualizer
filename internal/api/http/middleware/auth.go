package middleware

import (
	"context"
	"net/http"
	"strings"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/model"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizer"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizercontext"
)

// RegisterUserTypes adds to request's context allowed user types.
func RegisterUserTypes(types ...model.UserType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, visualizercontext.UserTypesKey, model.UserTypes(types))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Auth searches user by token and adds his data to context.
func Auth(visualizer *visualizer.Visualizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			allowedUserTypes := visualizercontext.GetUserTypes(ctx)
			token := getAuthToken(r)

			user, err := visualizer.GetUserByToken(ctx, token, allowedUserTypes...)
			if err != nil {
				log.
					FromContext(ctx).
					WithStr("token", token).
					Debug("verify token")

				api.RespondError(ctx, w, err)
				return
			}

			ctx = context.WithValue(ctx, visualizercontext.UserKey, user)
			ctx = context.WithValue(ctx, visualizercontext.TokenKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CheckPermission checks if user has permission to the action.
func CheckPermission(visualizer *visualizer.Visualizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user := visualizercontext.GetUser(ctx)
			if user == nil {
				api.RespondError(ctx, w, ierr.ErrUnauthorized)
				return
			}

			filter := dataprovider.
				NewPermissionFilter().
				ByUserIDs(user.ID).
				ByActions(r.Method + ":" + r.URL.Path[7:])

			if err := visualizer.CheckPermission(ctx, filter); err != nil {
				log.
					FromContext(ctx).
					WithErr(err).
					WithFields(
						"userID", user.ID,
						"actions", filter.Actions,
					).
					Debug("check permission")

				api.RespondError(ctx, w, err)
				return
			}

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
