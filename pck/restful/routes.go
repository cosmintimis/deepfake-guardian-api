package restful

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *restfulApi) Routes() http.Handler {
	router := createRouter()

	router.Route("/api/health-check", func(r chi.Router) {
		r.Get("/v1/status", app.serverStatus)
	})

	router.Route("/api/media", func(r chi.Router) {
		r.Get("/v1/{id}", app.getMediaById)
		r.Get("/v1", app.getAllMedia)
		r.Delete("/v1/{id}", app.deleteMediaById)
		r.Post("/v1", app.addNewMedia)
		r.Put("/v1/{id}", app.updateMedia)
	})

	router.Handle("/ws", http.HandlerFunc(app.wsHandler))

	return router
}

func createRouter() chi.Router {
	router := chi.NewMux()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	router.Use((cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})))

	return router
}
