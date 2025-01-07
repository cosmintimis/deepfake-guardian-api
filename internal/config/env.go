package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	ENV_DEV  string = "dev"
	ENV_PROD string = "prod"
)

type Config struct {
	Env         string `required:"true"`
	Port        string `required:"true"`
	ServerUrl   string `required:"true" envconfig:"SERVER_URL"`
	DatabaseUrl string `required:"true" envconfig:"DATABASE_URL"`
}

var globalConfig Config

func loadEnv(logger *slog.Logger) {
	logger.Info("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	if os.Getenv("ENV") != ENV_PROD {
		loadEnv(logger)
	} else {
		logger.Info("Skip loading .env file")
	}

	err := envconfig.Process("", &globalConfig)
	if err != nil {
		return nil, err
	}
	return &globalConfig, nil
}

func GetConfig() *Config {
	return &globalConfig
}
