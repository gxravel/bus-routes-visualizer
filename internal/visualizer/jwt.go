package visualizer

import (
	"context"

	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
)

func (r *Visualizer) VerifyToken(ctx context.Context, token string) error {
	if token == "" {
		err := ierr.NewReason(ierr.ErrInvalidToken)
		return err
	}

	err := r.tokenManager.Verify(ctx, token)
	return err
}
