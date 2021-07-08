package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gxravel/bus-routes-visualizer/assets"
	mw "github.com/gxravel/bus-routes-visualizer/internal/api/http/middleware"
	"github.com/gxravel/bus-routes-visualizer/internal/config"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/model"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizer"

	"github.com/go-chi/chi"
)

type Server struct {
	*http.Server
	logger       log.Logger
	visualizer   *visualizer.Visualizer
	busroutesAPI string
}

func NewServer(
	cfg *config.Config,
	visualizer *visualizer.Visualizer,
	logger log.Logger,
) *Server {
	srv := &Server{
		Server: &http.Server{
			Addr:         cfg.API.Address,
			ReadTimeout:  cfg.API.ReadTimeout,
			WriteTimeout: cfg.API.WriteTimeout,
		},
		logger:       logger.WithStr("module", "api:http"),
		visualizer:   visualizer,
		busroutesAPI: cfg.API.BusRoutes,
	}

	r := chi.NewRouter()

	r.Use(mw.Logger(srv.logger))
	r.Use(mw.Recoverer)

	if cfg.API.ServeSwagger {
		registerSwagger(r)
	}

	r.Route("/internal", func(r chi.Router) {
		r.Get("/health", srv.getHealth)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.Route("/permissions", func(r chi.Router) {
				r.Use(
					mw.RegisterUserTypes(model.UserAdmin),
					mw.Auth(visualizer),
				)

				r.Get("/", srv.getPermissions)
			})

			r.Route("/graphs", func(r chi.Router) {
				r.Use(
					mw.RegisterUserTypes(
						model.UserAdmin,
						model.UserService,
					),
					mw.Auth(visualizer),
					mw.CheckPermission(visualizer),
				)

				r.Get("/", srv.getGraph)
			})
		})
	})

	srv.Handler = r

	return srv
}

func registerSwagger(r *chi.Mux) {
	r.HandleFunc("/internal/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/internal/swagger/", http.StatusFound)
	})

	swaggerHandler := http.StripPrefix("/internal/", http.FileServer(http.FS(assets.SwaggerFiles)))
	r.Get("/internal/swagger/*", swaggerHandler.ServeHTTP)
}

// nolint
func (s *Server) processRequest(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		s.logger.WithErr(err).Error("decode data")
		return err
	}
	return nil
}
