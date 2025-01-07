package restful

import (
	"log/slog"

	"github.com/cosmintimis/deepfake-guardian-api/pck/healthcheck"
)

type restfulApi struct {
	logger      *slog.Logger
	healthcheck healthcheck.Service
}

func New(logger *slog.Logger, healthcheck healthcheck.Service) *restfulApi {
	return &restfulApi{
		logger:      logger,
		healthcheck: healthcheck,
	}
}
