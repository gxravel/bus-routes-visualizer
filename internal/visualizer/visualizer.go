package visualizer

import (
	"github.com/gxravel/bus-routes-visualizer/internal/config"
	"github.com/gxravel/bus-routes-visualizer/internal/database"
	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	"github.com/gxravel/bus-routes-visualizer/internal/jwt"
	"github.com/gxravel/bus-routes-visualizer/internal/logger"
)

type Visualizer struct {
	config       *config.Config
	db           *database.Client
	logger       logger.Logger
	txer         dataprovider.Txer
	tokenManager jwt.Manager
}

func New(
	config *config.Config,
	db *database.Client,
	logger logger.Logger,
	txer dataprovider.Txer,
	jwtManager jwt.Manager,
) *Visualizer {
	return &Visualizer{
		config:       config,
		db:           db,
		logger:       logger,
		txer:         txer,
		tokenManager: jwtManager,
	}
}
