package treatment

import (
	"vet-clinic-api/config"
	"vet-clinic-api/pkg/authentication"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) chi.Router {

	// Init router
	treatmentConfig := New(configuration)
	router := chi.NewRouter()

	// Routes protected by authentication
	router.Group(func(router chi.Router) {
		router.Use(authentication.AuthMiddleware(treatmentConfig.JWTSecret))

		router.Get("/", treatmentConfig.GetAllHandler)
		router.Get("/{id}", treatmentConfig.GetByIdHandler)
		router.Get("/{id}/history", treatmentConfig.GetByVisitIdHandler)
		router.Post("/", treatmentConfig.PostHandler)
		router.Put("/{id}", treatmentConfig.UpdateHandler)
		router.Delete("/{id}", treatmentConfig.DeleteHandler)
	})

	return router
}
