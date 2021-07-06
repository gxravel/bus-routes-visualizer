package visualizer

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/gxravel/bus-routes-visualizer/internal/busroutesapi"
	v1 "github.com/gxravel/bus-routes-visualizer/internal/busroutesapi/v1"
	"github.com/gxravel/bus-routes-visualizer/internal/drawing"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
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

// DrawGraph returns path to the graph image.
func (r *Visualizer) DrawGraph(ctx context.Context, routes []*v1.RouteDetailed) (int64, []byte, error) {
	graphName := routes[0].City + "_" + routes[0].Bus

	path, err := drawing.DrawRoutes(graphName, routes)
	if err != nil {
		return 0, nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		log.
			FromContext(ctx).
			WithErr(err).
			Error("open file")
		return 0, nil, err
	}
	defer file.Close()

	buf := &bytes.Buffer{}

	size, err := io.Copy(buf, file)
	if err != nil {
		log.
			FromContext(ctx).
			WithErr(err).
			Error("write to buffer")
		return 0, nil, err
	}

	return size, buf.Bytes(), nil
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
