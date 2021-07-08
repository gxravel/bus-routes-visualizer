package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

type PermissionStore interface {
	WithTx(*Tx) PermissionStore
	GetByFilter(ctx context.Context, filter *PermissionFilter) (*model.Permission, error)
	GetListByFilter(ctx context.Context, filter *PermissionFilter) ([]*model.Permission, error)
	Add(ctx context.Context, permissions ...*model.Permission) error
	Update(ctx context.Context, permission ...*model.Permission) ([]int64, error)
	Delete(ctx context.Context, filter *PermissionFilter) error
}

type PermissionFilter struct {
	UserIDs []int64
	Actions []string
}

func NewPermissionFilter() *PermissionFilter {
	return &PermissionFilter{}
}

// ByUserIDs filters by permission.user_id.
func (f *PermissionFilter) ByUserIDs(userIDs ...int64) *PermissionFilter {
	f.UserIDs = userIDs
	return f
}

// ByActions filters by permission.actions.
func (f *PermissionFilter) ByActions(actions ...string) *PermissionFilter {
	f.Actions = actions
	return f
}
