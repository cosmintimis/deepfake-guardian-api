package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/cosmintimis/deepfake-guardian-api/internal/config"
	"github.com/cosmintimis/deepfake-guardian-api/pck/healthcheck"
	"github.com/cosmintimis/deepfake-guardian-api/pck/postgresql"
	"github.com/cosmintimis/deepfake-guardian-api/pck/restful"
	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	logger.Info("Starting server")
	logger.Info("Loading configuration")
	// load .env file from given path
	// we keep it empty it will load .env from current directory
	config, configError := config.LoadConfig(logger)
	if configError != nil {
		log.Fatal(configError)
	}
	logger.Info("Configuration loaded", "env", config.Env)

	newConn, dbError := postgresql.InitDB()
	if dbError != nil {
		log.Fatal(dbError)
	}
	defer newConn.Close(context.Background())

	healthcheck := healthcheck.New()

	restfulApi := restful.New(logger, healthcheck)
	router := restfulApi.Routes()

	port := config.Port
	server := &http.Server{
		Handler: router,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}

	logger.Info("Server started on port " + port)
	log.Fatal(server.ListenAndServe())
}
