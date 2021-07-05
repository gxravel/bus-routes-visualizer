package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	busroutesapi "github.com/gxravel/bus-routes-visualizer/internal/busroutesapi"
)

func (s *Server) getGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	url, err := api.ParseGraphsRequest(r, s.busroutesAPI+busroutesapi.RouteForBuses)
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

	path, err := s.visualizer.DrawGraph(routes)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondPNG(ctx, w, path)
}
