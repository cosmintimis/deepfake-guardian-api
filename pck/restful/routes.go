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

	router.Route("/api/media", func(r chi.Router) {
		r.Get("/v1", app.getMediaById)
		r.Delete("/v1/{id}", app.deleteMediaById)
		r.Post("/v1", app.addNewMedia)
		r.Put("/v1/{id}", app.updateMedia)
		r.Get("/v1/all", app.getAllMedia)
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
