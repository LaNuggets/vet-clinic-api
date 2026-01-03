package cat

import (
	"vet-clinic-api/config"
	"vet-clinic-api/pkg/authentication"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) chi.Router {

	// Init router
	catConfig := New(configuration)
	router := chi.NewRouter()

	// Routes protected by authentication
	router.Group(func(router chi.Router) {
		router.Use(authentication.AuthMiddleware(catConfig.JWTSecret))

		router.Get("/{id}", catConfig.GetByIdHandler)
		router.Get("/{id}/history", catConfig.GetCatHistoryHandler)
		router.Get("/", catConfig.GetAllHandler)
		router.Post("/", catConfig.PostHandler)
		router.Put("/{id}", catConfig.UpdateHandler)
		router.Delete("/{id}", catConfig.DeleteHandler)
	})

	return router
}
