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

	routes, err := busroutesapi.GetRoutesDetailed(ctx, s.busroutesAPI, url)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if routes == nil {
		api.RespondEmptyItems(ctx, w)
		return
	}

	api.RespondDataOK(ctx, w, routes)
}
