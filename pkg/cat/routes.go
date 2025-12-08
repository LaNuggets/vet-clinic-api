package cat

import (
	"vet-clinic-api/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) chi.Router {

	// Init router
	catConfig := New(configuration)
	router := chi.NewRouter()

	// Routes definition
	router.Get("/{id}", catConfig.GetByIdHandler)
	router.Get("/{id}/history", catConfig.GetCatHistoryHandler)
	router.Get("/", catConfig.GetAllHandler)
	router.Post("/", catConfig.PostHandler)
	router.Put("/{id}", catConfig.UpdateHandler)
	router.Delete("/{id}", catConfig.DeleteHandler)

	return router
}
