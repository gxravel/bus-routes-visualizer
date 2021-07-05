package visualizercontext

import (
	"context"
)

type ctxKey string

const (
	TokenKey ctxKey = "token"
)

// GetToken returns request token, if auth is successful.
func GetToken(ctx context.Context) string {
	if ctx != nil {
		if val, ok := ctx.Value(TokenKey).(string); ok {
			return val
		}
	}

	return ""
}
