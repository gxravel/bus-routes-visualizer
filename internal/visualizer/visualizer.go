package visualizer

import (
	"github.com/gxravel/bus-routes-visualizer/internal/config"
	"github.com/gxravel/bus-routes-visualizer/internal/database"
	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	"github.com/gxravel/bus-routes-visualizer/internal/jwt"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/service"
)

type Visualizer struct {
	config          *config.Config
	db              *database.Client
	logger          log.Logger
	txer            dataprovider.Txer
	routeStore      dataprovider.RouteStore
	routePointStore dataprovider.RoutePointStore
	tokenManager    jwt.Manager
	busroutes       service.BusRoutes
}

func New(
	config *config.Config,
	db *database.Client,
	logger log.Logger,
	txer dataprovider.Txer,
	routeStore dataprovider.RouteStore,
	routePointStore dataprovider.RoutePointStore,
	jwtManager jwt.Manager,
	busroutes service.BusRoutes,
) *Visualizer {
	return &Visualizer{
		config:          config,
		db:              db,
		logger:          logger,
		txer:            txer,
		tokenManager:    jwtManager,
		routeStore:      routeStore,
		routePointStore: routePointStore,
		busroutes:       busroutes,
	}
}
