package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/marban004/factory_games_organizer/handler"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (a *AppDispatcher) loadRoutes() {
	usersHandler := &handler.Dispatcher{
		Secret:                           a.secret,
		UsersMicroservicesAddresses:      a.usersMicroservicesAddresses,
		CrudMicroservicesAddresses:       a.crudMicroservicesAddresses,
		CalculatorMicroservicesAddresses: a.calculatorMicroservicesAddresses,
		NextUsers:                        0,
		NextCrud:                         0,
		Nextcalculator:                   0,
	}
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Get("/health", usersHandler.Health)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("https://localhost:%d/swagger/doc.json", a.config.ServerPort)), //The url pointing to API definition
	))
	a.router = router
}
