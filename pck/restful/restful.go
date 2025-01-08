package restful

import (
	"log/slog"

	"github.com/cosmintimis/deepfake-guardian-api/pck/business/repositories"
	"github.com/cosmintimis/deepfake-guardian-api/pck/healthcheck"
	"github.com/cosmintimis/deepfake-guardian-api/pck/postgresql"
)

type restfulApi struct {
	logger          *slog.Logger
	healthcheck     healthcheck.Service
	mediaRepository repositories.MediaRepository
}

func New(logger *slog.Logger, healthcheck healthcheck.Service) *restfulApi {
	return &restfulApi{
		logger:          logger,
		healthcheck:     healthcheck,
		mediaRepository: postgresql.NewMediaRepository(logger),
	}
}
