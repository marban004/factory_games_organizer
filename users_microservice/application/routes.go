package application

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/marban004/factory_games_organizer/handler"
	"github.com/marban004/factory_games_organizer/microservice_logic_users/repository/user"
)

func (a *AppCrud) loadRoutes() {
	usersHandler := &handler.Users{
		UserRepo: &user.MySQLRepo{DB: a.db},
		Secret:   a.secret,
	}
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Post("/", usersHandler.CreateUser)
	router.Put("/", usersHandler.UpdateUser)
	a.router = router
}
