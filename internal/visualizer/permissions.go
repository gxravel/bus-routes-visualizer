package visualizer

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

func (r *Visualizer) GetPermissions(ctx context.Context, filter *dataprovider.PermissionFilter) ([]*httpv1.Permission, error) {
	dbPermissions, err := r.permissionStore.GetListByFilter(ctx, filter)
	if err != nil {
		log.
			FromContext(ctx).
			WithErr(err).
			Error("get permissions")

		return nil, err
	}

	return toV1Permissions(dbPermissions...), nil
}

func toV1Permissions(dbPermissions ...*model.Permission) []*httpv1.Permission {
	permissions := make([]*httpv1.Permission, 0, len(dbPermissions))
	for _, p := range dbPermissions {
		actions, ok := p.Actions["actions"]
		if !ok {
			continue
		}
		permissions = append(permissions, &httpv1.Permission{
			UserID:  p.UserID,
			Actions: actions,
		})
	}

	return permissions
}
