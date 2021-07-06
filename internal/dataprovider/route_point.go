package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

type RoutePointStore interface {
	WithTx(*Tx) RoutePointStore
	Add(ctx context.Context, points ...*model.RoutePoint) error
	Update(ctx context.Context, point *model.RoutePoint) error
	Delete(ctx context.Context, filter *RoutePointFilter) error
}

type RoutePointFilter struct {
	IDs []int64
}

func NewRoutePointFilter() *RoutePointFilter {
	return &RoutePointFilter{}
}

// ByIDs filters by route_point.id
func (f *RoutePointFilter) ByIDs(ids ...int64) *RoutePointFilter {
	f.IDs = ids
	return f
}
