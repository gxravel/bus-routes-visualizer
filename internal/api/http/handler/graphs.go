package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	service "github.com/gxravel/bus-routes-visualizer/internal/service/http"
)

func (s *Server) getGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	url, err := api.ParseGraphsRequest(r, s.busroutesAPI+service.RouteForBuses)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	routes, err := s.visualizer.GetRoutesDetailed(ctx, url)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if routes == nil {
		api.RespondEmptyItems(ctx, w)
		return
	}

	size, image, err := s.visualizer.GetRoutesGraph(ctx, routes)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondImageOK(ctx, w, size, image)
}
