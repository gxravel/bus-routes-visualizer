package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gxravel/bus-routes-visualizer/assets"
	mw "github.com/gxravel/bus-routes-visualizer/internal/api/http/middleware"
	"github.com/gxravel/bus-routes-visualizer/internal/config"
	"github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizer"

	"github.com/go-chi/chi"
)

type Server struct {
	*http.Server
	logger     logger.Logger
	visualizer *visualizer.Visualizer
}

func NewServer(
	cfg *config.Config,
	visualizer *visualizer.Visualizer,
	logger logger.Logger,
) *Server {
	srv := &Server{
		Server: &http.Server{
			Addr:         cfg.API.Address,
			ReadTimeout:  cfg.API.ReadTimeout,
			WriteTimeout: cfg.API.WriteTimeout,
		},
		logger:     logger.WithStr("module", "api:http"),
		visualizer: visualizer,
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

func (s *Server) processRequest(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		s.logger.WithErr(err).Error("decoding data")
		return err
	}
	return nil
}
