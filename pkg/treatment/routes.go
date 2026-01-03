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

		// Routes protected by authentication and accessible by admin only
		router.With(authentication.RoleMiddleware("admin")).Group(func(r chi.Router) {
			r.Post("/", treatmentConfig.PostHandler)
			r.Put("/{id}", treatmentConfig.UpdateHandler)
			r.Delete("/{id}", treatmentConfig.DeleteHandler)
		})
	})

	return router
}
