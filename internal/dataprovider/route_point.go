package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

type RoutePointStore interface {
	WithTx(*Tx) RoutePointStore
	Add(ctx context.Context, points ...*model.RoutePoint) error
	Update(ctx context.Context, point *model.RoutePoint) error
	Delete(ctx context.Context, id int64) error
}
