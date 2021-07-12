package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
)

func (s *Server) getHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := s.visualizer.IsHealthy(ctx); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondEmpty(w)
}
