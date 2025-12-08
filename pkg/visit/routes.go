package visit

import (
	"vet-clinic-api/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) chi.Router {

	// Init Router
	visitConfig := New(configuration)
	router := chi.NewRouter()

	// Routes definition
	router.Get("/", visitConfig.GetAlldHandler)
	router.Get("/{id}", visitConfig.GetByIdHandler)
	router.Post("/", visitConfig.PostHandler)
	router.Put("/{id}", visitConfig.UpdateHandler)
	router.Delete("/{id}", visitConfig.DeleteHandler)

	return router
}
