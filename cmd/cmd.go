package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/gxravel/bus-routes-visualizer/internal/api/http/handler"
	"github.com/gxravel/bus-routes-visualizer/internal/config"
	"github.com/gxravel/bus-routes-visualizer/internal/database"
	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider/mysql"
	"github.com/gxravel/bus-routes-visualizer/internal/jwt"
	"github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/storage"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizer"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	defaultLog := logger.Default()

	var configPath = flag.String("config", "./config.example.json", "path to configuration file")

	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		defaultLog.WithErr(err).Fatal("can't create config")
	}

	log, err := logger.New(
		cfg.Log.Level,
		logger.DefaultOutput)
	if err != nil {
		defaultLog.WithErr(err).Error("can not init log with specified params, defaults are used")
		log = defaultLog
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.CtxWithLogger(ctx, log)

	db, err := database.NewClient(*cfg, log)
	if err != nil {
		log.WithErr(err).Fatal("construct database")
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.WithErr(err).Error("close database connection")
		}
	}()

	if err := db.Migrate(); err != nil {
		log.WithErr(err).Fatal("can't migrate the db")
	}

	storage, err := storage.NewClient(*cfg)
	if err != nil {
		log.WithErr(err).Fatal("connect to storage client")
	}

	txer := mysql.NewTxManager(db)
	routeStore := mysql.NewRouteStore(db, txer)
	routePointStore := mysql.NewRoutePointStore(db, txer)

	jwtManager := jwt.New(storage, *cfg)

	visualizer := visualizer.New(
		cfg,
		db,
		log,
		txer,
		routeStore,
		routePointStore,
		jwtManager,
	)

	apiServer := handler.NewServer(
		cfg,
		visualizer,
		log,
	)

	var (
		shutdown     = make(chan os.Signal, 1)
		serverErrors = make(chan error, 1)
	)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		serverErrors <- apiServer.ListenAndServe()
	}()

	log.Info("started")

	defer log.Info("stopped")

	select {
	case err = <-serverErrors:
		log.WithErr(err).Error("api server stopped")

	case sig := <-shutdown:
		log.WithField("signal", sig.String()).Info("gracefully shutdown application")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err = apiServer.Shutdown(ctx); err != nil {
			log.WithErr(err).Error("api server shutdown error")
			err = apiServer.Close()
		}

		if err != nil {
			log.WithErr(err).Error("could not stopped api server gracefully")
		}
	}

}
