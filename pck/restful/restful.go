package restful

import (
	"log/slog"
	"sync"

	"github.com/cosmintimis/deepfake-guardian-api/pck/business/repositories"
	"github.com/cosmintimis/deepfake-guardian-api/pck/healthcheck"
	"github.com/cosmintimis/deepfake-guardian-api/pck/postgresql"
	"github.com/gorilla/websocket"
)

type restfulApi struct {
	logger          *slog.Logger
	healthcheck     healthcheck.Service
	mediaRepository repositories.MediaRepository
	connections     map[*websocket.Conn]struct{}
	connLock        sync.Mutex
}

func New(logger *slog.Logger, healthcheck healthcheck.Service) *restfulApi {
	return &restfulApi{
		logger:          logger,
		healthcheck:     healthcheck,
		mediaRepository: postgresql.NewMediaRepository(logger),
		connections:     make(map[*websocket.Conn]struct{}),
		connLock:        sync.Mutex{},
	}
}
