package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
)

func (s *Server) getPermissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := api.ParsePermissionsFilter(r)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	permissions, err := s.visualizer.GetPermissions(ctx, filter)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondDataOK(ctx, w, permissions)
}
