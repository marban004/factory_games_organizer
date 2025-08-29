package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/marban004/factory_games_organizer/handler"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (a *AppDispatcher) loadRoutes() {
	dispatcherhandler := &handler.Dispatcher{
		UsersMicroservicesAddresses:      a.usersMicroservicesAddresses,
		CrudMicroservicesAddresses:       a.crudMicroservicesAddresses,
		CalculatorMicroservicesAddresses: a.calculatorMicroservicesAddresses,
	}
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(a.statTracker.ApiStatTracker)
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
	router.Get("/health", dispatcherhandler.Health)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("https://%s:%d/swagger/doc.json", a.config.Host, a.config.ServerPort)), //The url pointing to API definition
	))
	router.Route("/users", a.loadUserRoutes)
	router.Route("/crud", a.loadCrudRoutes)
	router.Route("/calculator", a.loadCalculatorRoutes)
	a.router = router
}

func (a *AppDispatcher) loadUserRoutes(router chi.Router) {
	dispatcherHandlerUsers := handler.DispatcherUsers{
		CommonHandlerFunctions: handler.CommonHandlerFunctions{
			Secret:           a.secret,
			NextMicroservice: 0,
			Client: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
		UsersMicroservicesAddresses: a.usersMicroservicesAddresses,
	}
	router.Post("/login", dispatcherHandlerUsers.LoginUser)
	router.Post("/", dispatcherHandlerUsers.CreateUser)
	router.Put("/", dispatcherHandlerUsers.UpdateUser)
	router.Delete("/", dispatcherHandlerUsers.DeleteUser)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("https://%s/swagger/doc.json", dispatcherHandlerUsers.UsersMicroservicesAddresses[0])), //The url pointing to API definition
	))
}

func (a *AppDispatcher) loadCrudRoutes(router chi.Router) {
	dispatcherHandlerCrud := handler.DispatcherCrud{
		CommonHandlerFunctions: handler.CommonHandlerFunctions{
			Secret:           a.secret,
			NextMicroservice: 0,
			Client: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
		CrudMicroservicesAddresses: a.crudMicroservicesAddresses,
	}
	router.Get("/selectbyid", dispatcherHandlerCrud.SelectByID)
	router.Get("/select", dispatcherHandlerCrud.Select)
	router.Post("/", dispatcherHandlerCrud.Insert)
	router.Put("/", dispatcherHandlerCrud.Update)
	router.Delete("/", dispatcherHandlerCrud.Delete)
	router.Delete("/user", dispatcherHandlerCrud.DeleteByUser)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("https://%s/swagger/doc.json", dispatcherHandlerCrud.CrudMicroservicesAddresses[0])), //The url pointing to API definition
	))
}

func (a *AppDispatcher) loadCalculatorRoutes(router chi.Router) {
	dispatcherHandlerCalculator := handler.DispatcherCalculator{
		CommonHandlerFunctions: handler.CommonHandlerFunctions{
			Secret:           a.secret,
			NextMicroservice: 0,
			Client: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
		CalculatorMicroservicesAddresses: a.calculatorMicroservicesAddresses,
	}
	router.Get("/calculate", dispatcherHandlerCalculator.Calculate)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("https://%s/swagger/doc.json", dispatcherHandlerCalculator.CalculatorMicroservicesAddresses[0])), //The url pointing to API definition
	))
}
