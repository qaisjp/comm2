package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/multitheftauto/community/pkg/api"
	"github.com/multitheftauto/community/pkg/config"
	"github.com/multitheftauto/community/pkg/database"

	"github.com/jmoiron/sqlx"
	"github.com/koding/multiconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error

	m := multiconfig.NewWithPath(os.Getenv("config"))
	cfg := &config.Config{}
	m.MustLoad(cfg)

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger := logrus.StandardLogger()
	logger.Level = logLevel

	logger.WithFields(logrus.Fields{
		"module": "init",
	}).Info("Starting up the application")

	// Initialize the database
	var db *sqlx.DB

	db, err = database.NewPostgres(cfg.Postgres)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"error":  err.Error(),
			"cstr":   cfg.Postgres.ConnectionString,
		}).Fatal("Unable to connect to the Postgres server")
		return
	}

	logger.WithFields(logrus.Fields{
		"module": "init",
		"cstr":   cfg.Postgres.ConnectionString,
	}).Info("Connected to a Postgres server")

	api := api.NewAPI(
		cfg,
		logger,
		db,
	)

	go func() {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"bind":   cfg.Address,
		}).Info("Starting the API server")

		if err := api.Start(); err != nil {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"error":  err.Error(),
			}).Fatal("API server failed")
		}
	}()

	// Create a new signal receiver
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Watch for a signal
	<-sc

	// ugly thing to stop ^C from killing alignment
	logger.Out.Write([]byte("\r\n"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := api.Shutdown(ctx); err != nil {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"error":  err.Error(),
		}).Fatal("Failed to close the API server")
	}

	logger.WithFields(logrus.Fields{
		"module": "init",
	}).Info("mtacommunity-api has shut down.")
}
