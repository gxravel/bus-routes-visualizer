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

// SetPermissions creates permissions for users that don't have them,
// and updates for the ones who do.
func (r *Visualizer) SetPermissions(ctx context.Context, permissions []*httpv1.Permission) error {
	dbPermissions := toDBPermissions(permissions...)

	affectedIDs, err := r.permissionStore.Update(ctx, dbPermissions...)
	if err != nil {
		log.
			FromContext(ctx).
			WithErr(err).
			Error("failed to update permissions")

		return err
	}

	skipIDs := make(map[int64]struct{}, len(affectedIDs))
	for _, id := range affectedIDs {
		skipIDs[id] = struct{}{}
	}

	for i, p := range dbPermissions {
		if _, should := skipIDs[p.UserID]; should {
			dbPermissions[i] = dbPermissions[len(dbPermissions)-1]
			dbPermissions = dbPermissions[:len(dbPermissions)-1]
		}
	}

	if len(dbPermissions) == 0 {
		return nil
	}

	if err := r.permissionStore.Add(ctx, dbPermissions...); err != nil {
		log.
			FromContext(ctx).
			WithErr(err).
			Error("failed to add permissions")

		if derr := ierr.CheckDuplicate(err, "user_id"); derr != nil {
			return derr
		}

		return err
	}

	return nil
}

// DeletePermissions deletes users permissions for actions (by filter).
func (r *Visualizer) DeletePermissions(ctx context.Context, filter *dataprovider.PermissionFilter) error {
	return r.permissionStore.Delete(ctx, filter)
}

// CheckPermission checks if user has permission for actions (by filter).
func (r *Visualizer) CheckPermission(ctx context.Context, filter *dataprovider.PermissionFilter) error {
	dbPermission, err := r.permissionStore.GetByFilter(ctx, filter)
	if err != nil {
		return err
	}

	if dbPermission == nil {
		return ierr.ErrPermissionDenied
	}

	return nil
}

func toDBPermissions(permissions ...*httpv1.Permission) []*model.Permission {
	dbPermissions := make([]*model.Permission, 0, len(permissions))
	for _, p := range permissions {
		dbPermissions = append(dbPermissions, &model.Permission{
			UserID: p.UserID,
			Actions: model.JSON{
				"actions": p.Actions,
			},
		})
	}

	return dbPermissions
}

func toV1Permissions(dbPermissions ...*model.Permission) []*httpv1.Permission {
	permissions := make([]*httpv1.Permission, 0, len(dbPermissions))
	for _, p := range dbPermissions {
		actions, should := p.Actions["actions"]
		if !should {
			continue
		}
		permissions = append(permissions, &httpv1.Permission{
			UserID:  p.UserID,
			Actions: actions,
		})
	}

	return permissions
}
