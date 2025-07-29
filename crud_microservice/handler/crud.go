package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/marban004/factory_games_organizer/microservice_logic_crud/model"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine"
	machinerecipe "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine_recipe"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe"
	recipeinput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_input"
	recipeoutput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_output"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/resource"
)

type JSONData struct {
	MachinesList        []model.MachineInfo
	ResourcesList       []model.ResourceInfo
	RecipesList         []model.RecipeInfo
	RecipesInputsList   []model.RecipeInputOutputInfo
	RecipesOutputsList  []model.RecipeInputOutputInfo
	MachinesRecipesList []model.MachinesRecipesInfo
	JWT                 string
}

type CRUD struct {
	MachineRepo       *machine.MySQLRepo
	ResourceRepo      *resource.MySQLRepo
	RecipeRepo        *recipe.MySQLRepo
	RecipeinputRepo   *recipeinput.MySQLRepo
	RecipeoutputRepo  *recipeoutput.MySQLRepo
	MachineRecipeRepo *machinerecipe.MySQLRepo
	Secret            string
}

func (h *CRUD) SelectByID(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//jwt = token with dispatcher server secret key, id of user who received the token and issue date of the token, not optional
	//machines_id = id of the record from the machines table to be retreived, optional, if multiple values are associated with identifier, records for id each will be retrieved
	//resources_id = id of the record from the resources table to be retreived, optional, if multiple values are associated with identifier, records for id each will be retrieved
	//recipes_ids = id of the record from the recipes table to be retreived, optional, if multiple values are associated with identifier, records for id each will be retrieved
	//recipes_inputs_id = id of the record from the recipes_inputs table to be retreived, optional, if multiple values are associated with identifier, records for id each will be retrieved
	//recipes_outputs_id = id of the record from the recipes_inputs table to be retreived, optional, if multiple values are associated with identifier, records for id each will be retrieved
	//machines_recipes_id = id of the record from the machines_recipes table to be retreived, optional, this table holds records of which machine can use which recipe, if multiple values are associated with identifier, records for id each will be retrieved
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, userId := h.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	returnData := JSONData{}
	machinesIds := r.URL.Query()["machines_id"]
	resourcesIds := r.URL.Query()["resources_id"]
	recipesIds := r.URL.Query()["recipes_id"]
	recipesInputsIds := r.URL.Query()["recipes_inputs_id"]
	recipesOutputsIds := r.URL.Query()["recipes_outputs_id"]
	machinesRecipesIds := r.URL.Query()["machines_recipes_id"]
	if machinesIds != nil {
		result, err := h.MachineRepo.SelectMachinesById(r.Context(), h.convertArrToInt(machinesIds), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
			return
		}
		returnData.MachinesList = result
	}
	if resourcesIds != nil {
		result, err := h.ResourceRepo.SelectResourcesById(r.Context(), h.convertArrToInt(resourcesIds), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
			return
		}
		returnData.ResourcesList = result
	}
	if recipesIds != nil {
		result, err := h.RecipeRepo.SelectRecipesById(r.Context(), h.convertArrToInt(recipesIds), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
			return
		}
		returnData.RecipesList = result
	}
	if recipesInputsIds != nil {
		result, err := h.RecipeinputRepo.SelectRecipesInputsById(r.Context(), h.convertArrToInt(recipesInputsIds), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
			return
		}
		returnData.RecipesInputsList = result
	}
	if recipesOutputsIds != nil {
		result, err := h.RecipeoutputRepo.SelectRecipesOutputsById(r.Context(), h.convertArrToInt(recipesOutputsIds), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
			return
		}
		returnData.RecipesOutputsList = result
	}
	if machinesRecipesIds != nil {
		result, err := h.MachineRecipeRepo.SelectMachinesRecipesById(r.Context(), h.convertArrToInt(machinesRecipesIds), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
			return
		}
		returnData.MachinesRecipesList = result
	}
	byteJSONRepresentation, err := json.Marshal(returnData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not generate json representation of data, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

func (h *CRUD) Select(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//userid = user with whose data we want to generate production tree, not optional
	//resource = resource for which we want to generate production tree, not optional
	//rate = target production rate per second for specified resource, not optional
	//alt_recipes = recipes that can be used besides default recipes, optional
	//alt_machines = machines that can be used besides default machines, optional
}

func (h *CRUD) Insert(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//userid = user with whose data we want to generate production tree, not optional
	//resource = resource for which we want to generate production tree, not optional
	//rate = target production rate per second for specified resource, not optional
	//alt_recipes = recipes that can be used besides default recipes, optional
	//alt_machines = machines that can be used besides default machines, optional
}

func (h *CRUD) Update(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//userid = user with whose data we want to generate production tree, not optional
	//resource = resource for which we want to generate production tree, not optional
	//rate = target production rate per second for specified resource, not optional
	//alt_recipes = recipes that can be used besides default recipes, optional
	//alt_machines = machines that can be used besides default machines, optional
}

func (h *CRUD) Delete(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//userid = user with whose data we want to generate production tree, not optional
	//resource = resource for which we want to generate production tree, not optional
	//rate = target production rate per second for specified resource, not optional
	//alt_recipes = recipes that can be used besides default recipes, optional
	//alt_machines = machines that can be used besides default machines, optional
}

// todo: implement verification of jwt
func (h *CRUD) verifyJWT(jwt string) (bool, int) {
	if len(jwt) > 0 {
		return true, 1
	}
	return false, 0
}

func (h *CRUD) convertArrToInt(input []string) []int {
	result := []int{}
	for _, value := range input {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		result = append(result, intValue)
	}
	return result
}
