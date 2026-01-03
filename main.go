package main

import (
	"log"
	"net/http"
	"vet-clinic-api/config"
	"vet-clinic-api/pkg/cat"
	"vet-clinic-api/pkg/treatment"
	"vet-clinic-api/pkg/user"
	"vet-clinic-api/pkg/visit"

	_ "vet-clinic-api/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Configuration of CORS middleware
	router.Use(cors.Handler(cors.Options{

		// Origins allowrd
		AllowedOrigins: []string{"https://*", "http://*"},

		// Allowed Methods for the requests
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},

		// Allowed headers
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},

		// Exposed Headers
		ExposedHeaders: []string{"Link"},

		// Allow Credentials
		AllowCredentials: true,

		// cache request time
		MaxAge: 300,
	}))

	// Set up content type
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// Route Set up
	router.Mount("/api/v1/vet/cats", cat.Routes(configuration))
	router.Mount("/api/v1/vet/treatments", treatment.Routes(configuration))
	router.Mount("/api/v1/vet/visits", visit.Routes(configuration))
	router.Mount("/api/v1/vet/users", user.Routes(configuration))

	// Load static file for swagger
	router.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	// Swagger documentation route
	router.Handle("/swagger/*", httpSwagger.WrapHandler)
	return router
}

// @title           Veterinarian API
// @version         1.0
// @description    	This is an API for managing a veterinary clinic. You can register cats, consultations and treatments.
// @host            localhost:8081
// @BasePath        /api/v1/vet
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Init configuration
	configuration, err := config.New()
	if err != nil {
		log.Panicln("Configuration error:", err)
	}

	// Init Routes
	router := Routes(configuration)

	// Lunch the server
	log.Println("Serving on :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
