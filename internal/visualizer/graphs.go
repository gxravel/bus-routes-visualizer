package visualizer

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/busroutesapi"
	v1 "github.com/gxravel/bus-routes-visualizer/internal/busroutesapi/v1"
	"github.com/gxravel/bus-routes-visualizer/internal/drawing"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

// GetRoutesDetailed fetches detailed routes and saves them in database.
func (r *Visualizer) GetRoutesDetailed(ctx context.Context, url string) ([]*v1.RouteDetailed, error) {
	routes, err := busroutesapi.GetRoutesDetailed(ctx, r.config.API.BusRoutes, url)
	if err != nil {
		return nil, err
	}
	if routes == nil {
		return nil, nil
	}

	err = r.routeStore.Add(ctx, toDBRoutes(routes)...)
	if err != nil && ierr.CheckDuplicate(err, "route") == nil {
		return nil, err
	}
	return routes, nil
}

func toDBRoutes(routes []*v1.RouteDetailed) []*model.Route {
	var dbRoutes = make([]*model.Route, 0, len(routes))
	for i, r := range routes {
		dbRoutes = append(dbRoutes, &model.Route{
			Bus:    r.Bus,
			City:   r.City,
			Points: make([]*model.RoutePoint, 0, len(r.Points)),
		})

		for _, p := range r.Points {
			dbRoutes[i].Points = append(dbRoutes[i].Points, &model.RoutePoint{
				Step:    p.Step,
				Address: p.Address,
			})
		}
	}

	return dbRoutes
}
