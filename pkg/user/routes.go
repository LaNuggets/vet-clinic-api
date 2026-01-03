package user

import (
	"vet-clinic-api/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) chi.Router {

	// Init router
	userConfig := New(configuration)
	router := chi.NewRouter()

	// Routes definition
	router.Post("/login", userConfig.LoginHandler)
	router.Post("/refresh", userConfig.RefreshHandler)
	router.Get("/{id}", userConfig.GetByIdHandler)
	router.Get("/", userConfig.GetAllHandler)
	router.Post("/", userConfig.PostHandler)
	router.Put("/{id}", userConfig.UpdateHandler)
	router.Delete("/{id}", userConfig.DeleteHandler)

	return router
}
