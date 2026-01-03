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

		// Routes protected by authentication and accessible by admin only
		router.With(authentication.RoleMiddleware("admin")).Group(func(r chi.Router) {
			r.Post("/", catConfig.PostHandler)
			r.Put("/{id}", catConfig.UpdateHandler)
			r.Delete("/{id}", catConfig.DeleteHandler)
		})
	})

	return router
}
