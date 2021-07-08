package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
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

func (s *Server) setPermissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	permissions := make([]*httpv1.Permission, 0)
	if err := s.processRequest(r, &permissions); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	s.logger.WithField("permissions", permissions[0]).Debug("decoded permissions")

	if err := s.visualizer.SetPermissions(ctx, permissions); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondNoContent(w)
}

func (s *Server) deletePermissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := api.ParsePermissionsFilter(r)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := s.visualizer.DeletePermissions(ctx, filter); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondNoContent(w)
}
