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

		// Routes protected by authentication and accessible by admin only
		router.With(authentication.RoleMiddleware("admin")).Group(func(r chi.Router) {
			r.Post("/", visitConfig.PostHandler)
			r.Put("/{id}", visitConfig.UpdateHandler)
			r.Delete("/{id}", visitConfig.DeleteHandler)
		})
	})

	return router
}
