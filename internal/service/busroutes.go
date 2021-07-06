package service

import (
	"context"

	v1 "github.com/gxravel/bus-routes-visualizer/internal/service/http/v1"
)

// BusRoutes defines busroutes api service.
type BusRoutes interface {
	GetRoutesDetailed(ctx context.Context, url string) ([]*v1.RouteDetailed, error)
}
