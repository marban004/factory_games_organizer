package application

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/marban004/factory_games_organizer/handler"
)

func (a *AppCalculator) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/calculate", a.loadCalculateRoutes)

	a.router = router
}

func (a *AppCalculator) loadCalculateRoutes(router chi.Router) {
	calculatorHandler := &handler.Calculator{
		DB: a.db,
	}
	router.Get("/", calculatorHandler.Calculate)
}
