package visit

import (
	"vet-clinic-api/config"
	"vet-clinic-api/pkg/authentication"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) chi.Router {

	// Init Router
	visitConfig := New(configuration)
	router := chi.NewRouter()

	// Routes protected by authentication
	router.Group(func(router chi.Router) {
		router.Use(authentication.AuthMiddleware(visitConfig.JWTSecret))

		router.Get("/", visitConfig.GetAlldHandler)
		router.Get("/{id}", visitConfig.GetByIdHandler)
		router.Post("/", visitConfig.PostHandler)
		router.Put("/{id}", visitConfig.UpdateHandler)
		router.Delete("/{id}", visitConfig.DeleteHandler)
	})

	return router
}
