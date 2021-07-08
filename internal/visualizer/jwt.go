package visualizer

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	"github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

func (r *Visualizer) GetUserByToken(ctx context.Context, token string, allowedTypes ...model.UserType) (*httpv1.User, error) {
	if token == "" {
		err := ierr.NewReason(ierr.ErrInvalidToken)
		return nil, err
	}

	jwtUser, err := r.tokenManager.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	user := &httpv1.User{
		ID:   jwtUser.ID,
		Type: jwtUser.Type,
	}

	types := model.UserTypes(allowedTypes)
	if !types.Exists(user.Type) {
		err := ierr.NewReason(ierr.ErrPermissionDenied)

		logger.
			FromContext(ctx).
			WithField("user", user).
			Warn(err.Error())

		return nil, err
	}

	return user, nil
}
