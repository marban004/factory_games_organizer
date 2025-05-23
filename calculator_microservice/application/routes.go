package application

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/marban004/factory_games_organizer/calculator_microservice/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/calculate", loadCalculateRoutes)

	return router
}

func loadCalculateRoutes(router chi.Router) {
	calculatorHandler := &handler.Calculator{}
	router.Get("/", calculatorHandler.Calculate)
}
