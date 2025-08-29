package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/marban004/factory_games_organizer/handler"
	"github.com/marban004/factory_games_organizer/microservice_logic_users/repository/user"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (a *AppUsers) loadRoutes() {
	usersHandler := &handler.Users{
		UserRepo: &user.MySQLRepo{DB: a.db},
		Secret:   a.secret,
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
	router.Post("/login", usersHandler.LoginUser)
	router.Post("/", usersHandler.CreateUser)
	router.Put("/", usersHandler.UpdateUser)
	router.Delete("/", usersHandler.DeleteUser)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("https://%s:%d/swagger/doc.json", a.config.Host, a.config.ServerPort)), //The url pointing to API definition
	))
	a.router = router
}
