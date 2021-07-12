package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
)

func (s *Server) getRoutesGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bus, err := api.ParseGraphsRequest(r)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	routes, err := s.visualizer.GetRoutesDetailed(ctx, bus)
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
