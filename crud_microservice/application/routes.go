package application

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/marban004/factory_games_organizer/handler"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine"
	machinerecipe "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine_recipe"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe"
	recipeinput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_input"
	recipeoutput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_output"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/resource"
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
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/selectbyid", crudHandler.SelectByID)
	router.Get("/select", crudHandler.Select)
	router.Post("/", crudHandler.Insert)
	router.Put("/", crudHandler.Update)
	router.Delete("/", crudHandler.Delete)
	router.Delete("/user", crudHandler.DeleteByUser)

	a.router = router
}
