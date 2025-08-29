package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/marban004/factory_games_organizer/handler"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine"
	machinerecipe "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine_recipe"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe"
	recipeinput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_input"
	recipeoutput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_output"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/resource"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *AppCrud) loadRoutes() {
	crudHandler := &handler.CRUD{
		MachineRepo:       &machine.MySQLRepo{DB: a.db},
		ResourceRepo:      &resource.MySQLRepo{DB: a.db},
		RecipeRepo:        &recipe.MySQLRepo{DB: a.db},
		RecipeinputRepo:   &recipeinput.MySQLRepo{DB: a.db},
		RecipeoutputRepo:  &recipeoutput.MySQLRepo{DB: a.db},
		MachineRecipeRepo: &machinerecipe.MySQLRepo{DB: a.db},
		Secret:            a.secret,
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
	router.Get("/health", crudHandler.Health)
	router.Get("/selectbyid", crudHandler.SelectByID)
	router.Get("/select", crudHandler.Select)
	router.Post("/", crudHandler.Insert)
	router.Put("/", crudHandler.Update)
	router.Delete("/", crudHandler.Delete)
	router.Delete("/user", crudHandler.DeleteByUser)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("https://%s:%d/swagger/doc.json", a.config.Host, a.config.ServerPort)), //The url pointing to API definition
	))

	a.router = router
}
