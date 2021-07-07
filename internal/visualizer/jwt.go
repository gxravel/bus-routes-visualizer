package visualizer

import (
	"context"

	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
)

func (r *Visualizer) VerifyToken(ctx context.Context, token string) (int64, error) {
	if token == "" {
		err := ierr.NewReason(ierr.ErrInvalidToken)
		return 0, err
	}

	return r.tokenManager.Verify(ctx, token)
}
