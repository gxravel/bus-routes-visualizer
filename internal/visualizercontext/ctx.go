package visualizercontext

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

type ctxKey string

const (
	UserTypesKey ctxKey = "user_types"
	TokenKey     ctxKey = "token"
	UserKey      ctxKey = "user"
)

// GetUserTypes returns registered user types.
func GetUserTypes(ctx context.Context) model.UserTypes {
	if ctx == nil {
		return nil
	}

	t, _ := ctx.Value(UserTypesKey).(model.UserTypes)

	return t
}

// GetToken returns request token, if auth is successful.
func GetToken(ctx context.Context) string {
	if ctx != nil {
		if val, ok := ctx.Value(TokenKey).(string); ok {
			return val
		}
	}

	return ""
}

// GetUser returns user, if auth is successful.
func GetUser(ctx context.Context) *httpv1.User {
	if ctx != nil {
		if val, ok := ctx.Value(UserKey).(*httpv1.User); ok {
			return val
		}
	}

	return nil
}
