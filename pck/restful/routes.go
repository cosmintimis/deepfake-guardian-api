package restful

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *restfulApi) Routes() http.Handler {
	router := createRouter()

	router.Route("/api/health-check", func(r chi.Router) {
		r.Get("/v1/status", app.serverStatus)
	})

	return router
}

func createRouter() chi.Router {
	router := chi.NewMux()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	return router
}
