package handler

import (
	"encoding/json"
	"errors"
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
	//recipes_id = id of the record from the recipes table to be retreived, optional, if multiple values are associated with identifier, records for id each will be retrieved
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
	//test url 127.0.0.1:3000/selectbyid?jwt=l&machines_id=1&machines_id=2&resources_id=1&resources_id=2&recipes_id=1&recipes_id=2&recipes_inputs_id=1&recipes_inputs_id=2&recipes_outputs_id=1&recipes_outputs_id=2&machines_recipes_id=1&machines_recipes_id=2
}

func (h *CRUD) Select(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//jwt = token with dispatcher server secret key, id of user who received the token and issue date of the token, not optional
	//machines_id_start = the id of first record to be retreived from machines table, optional, if absent retreives records from the first record in database for the user
	//machines_rows = the number of records to be retreived from machines table, optional, if absent retreives as many records as possible for the user
	//resources_id_start = the id of first record to be retreived from resources table, optional, if absent retreives records from the first record in database for the user
	//resources_rows = the number of records to be retreived from resources table, optional, if absent retreives as many records as possible for the user
	//recipes_id_start = the id of first record to be retreived from recipes table, optional, if absent retreives records from the first record in database for the user
	//recipes_rows = the number of records to be retreived from recipes table, optional, if absent retreives as many records as possible for the user
	//recipes_inputs_id_start = the id of first record to be retreived from recipes_inputs table, optional, if absent retreives records from the first record in database for the user
	//recipes_inputs_rows = the number of records to be retreived from recipes_inputs table, optional, if absent retreives as many records as possible for the user
	//recipes_outputs_id_start = the id of first record to be retreived from recipes_outputs table, optional, if absent retreives records from the first record in database for the user
	//recipes_outputs_rows = the number of records to be retreived from recipes_outputs table, optional, if absent retreives as many records as possible for the user
	//machines_resources_id_start = the id of first record to be retreived from machines_resources table, optional, if absent retreives records from the first record in database for the user
	//machines_resources_rows = the number of records to be retreived from machines_resources table, optional, if absent retreives as many records as possible for the user
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

	machinesIdStart, err := strconv.Atoi(r.URL.Query().Get("machines_id_start"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || machinesIdStart < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("machines_id_start should be a positive integer"))
		return
	}
	machinesIdStart = 0
	machinesRows, err := strconv.Atoi(r.URL.Query().Get("machines_rows"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || machinesRows < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("machines_rows should be a positive integer"))
		return
	}
	machinesRows = 0
	machinesResult, err := h.MachineRepo.SelectMachines(r.Context(), machinesIdStart, machinesRows, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
		return
	}
	returnData.MachinesList = machinesResult

	resourcesIdStart, err := strconv.Atoi(r.URL.Query().Get("resources_id_start"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || resourcesIdStart < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("resources_id_start should be a positive integer"))
		return
	}
	resourcesIdStart = 0
	resourcesRows, err := strconv.Atoi(r.URL.Query().Get("resources_rows"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || resourcesRows < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("resources_rows should be a positive integer"))
		return
	}
	resourcesRows = 0
	resourcesResult, err := h.ResourceRepo.SelectResources(r.Context(), resourcesIdStart, resourcesRows, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
		return
	}
	returnData.ResourcesList = resourcesResult

	recipesIdStart, err := strconv.Atoi(r.URL.Query().Get("recipes_id_start"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || recipesIdStart < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("recipes_id_start should be a positive integer"))
		return
	}
	recipesIdStart = 0
	recipesRows, err := strconv.Atoi(r.URL.Query().Get("recipes_rows"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || recipesRows < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("recipes_rows should be a positive integer"))
		return
	}
	recipesRows = 0
	recipesResult, err := h.RecipeRepo.SelectRecipes(r.Context(), recipesIdStart, recipesRows, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
		return
	}
	returnData.RecipesList = recipesResult

	recipesInputsIdStart, err := strconv.Atoi(r.URL.Query().Get("recipes_inputs_id_start"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || recipesInputsIdStart < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("recipes_inputs_id_start should be a positive integer"))
		return
	}
	recipesInputsIdStart = 0
	recipesInputsRows, err := strconv.Atoi(r.URL.Query().Get("recipes_inputs_rows"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || recipesInputsRows < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("recipes_inputs_rows should be a positive integer"))
		return
	}
	recipesInputsRows = 0
	recipesInputsResult, err := h.RecipeinputRepo.SelectRecipesInputs(r.Context(), recipesInputsIdStart, recipesInputsRows, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
		return
	}
	returnData.RecipesInputsList = recipesInputsResult

	recipesOutputsIdStart, err := strconv.Atoi(r.URL.Query().Get("recipes_outputs_id_start"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || recipesOutputsIdStart < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("recipes_outputs_id_start should be a positive integer"))
		return
	}
	recipesOutputsIdStart = 0
	recipesOutputsRows, err := strconv.Atoi(r.URL.Query().Get("recipes_outputs_rows"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || recipesOutputsRows < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("recipes_outputs_rows should be a positive integer"))
		return
	}
	recipesOutputsRows = 0
	recipesOutputsResult, err := h.RecipeoutputRepo.SelectRecipesOutputs(r.Context(), recipesOutputsIdStart, recipesOutputsRows, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
		return
	}
	returnData.RecipesOutputsList = recipesOutputsResult

	machinesRecipesIdStart, err := strconv.Atoi(r.URL.Query().Get("machines_recipes_id_start"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || machinesRecipesIdStart < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("machines_recipes_id_start should be a positive integer"))
		return
	}
	machinesRecipesIdStart = 0
	machinesRecipesRows, err := strconv.Atoi(r.URL.Query().Get("machines_recipes_rows"))
	if (err != nil && !errors.Is(err, strconv.ErrSyntax)) || machinesRecipesRows < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("machines_recipes_rows should be a positive integer"))
		return
	}
	machinesRecipesRows = 0
	machinesRecipesResult, err := h.MachineRecipeRepo.SelectMachinesRecipes(r.Context(), machinesRecipesIdStart, machinesRecipesRows, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not retrieve data from database, reason: %w", err).Error()))
		return
	}
	returnData.MachinesRecipesList = machinesRecipesResult

	byteJSONRepresentation, err := json.Marshal(returnData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not generate json representation of data, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
	//test url 127.0.0.1:3000/select?jwt=l
}

func (h *CRUD) Insert(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//jwt = token with dispatcher server secret key, id of user who received the token and issue date of the token, not optional
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
	inputData := JSONData{}
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("could not parse received body, reason: %w", err).Error()))
		return
	}
	var response struct {
		MachinesInserted        uint
		ResourcesInserted       uint
		RecipesInserted         uint
		RecipesInputsInserted   uint
		RecipesOutputsInserted  uint
		MachinesRecipesInserted uint
	}
	response.MachinesInserted = 0
	response.ResourcesInserted = 0
	response.RecipesInserted = 0
	response.RecipesInputsInserted = 0
	response.RecipesOutputsInserted = 0
	response.MachinesRecipesInserted = 0
	skipRows := false
	if inputData.MachinesList != nil {
		for _, entry := range inputData.MachinesList {
			entry.UsersId = uint(userId)
		}
		result, err := h.MachineRepo.InsertMachines(r.Context(), inputData.MachinesList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not insert requested machines data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.MachinesInserted = uint(noRows)
		}
	}
	if inputData.ResourcesList != nil {
		for _, entry := range inputData.ResourcesList {
			entry.UsersId = uint(userId)
		}
		result, err := h.ResourceRepo.InsertResources(r.Context(), inputData.ResourcesList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not insert requested resources data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.ResourcesInserted = uint(noRows)
		}
	}
	if inputData.RecipesList != nil {
		for _, entry := range inputData.RecipesList {
			entry.UsersId = uint(userId)
		}
		result, err := h.RecipeRepo.InsertRecipes(r.Context(), inputData.RecipesList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not insert requested recipes data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.RecipesInserted = uint(noRows)
		}
	}
	if inputData.RecipesInputsList != nil {
		for _, entry := range inputData.RecipesInputsList {
			entry.UsersId = uint(userId)
		}
		result, err := h.RecipeinputRepo.InsertRecipesInputs(r.Context(), inputData.RecipesInputsList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not insert requested recipes_inputs data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.RecipesInputsInserted = uint(noRows)
		}
	}
	if inputData.RecipesOutputsList != nil {
		for _, entry := range inputData.RecipesOutputsList {
			entry.UsersId = uint(userId)
		}
		result, err := h.RecipeoutputRepo.InsertRecipesOutputs(r.Context(), inputData.RecipesOutputsList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not insert requested recipes_outputs data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.RecipesOutputsInserted = uint(noRows)
		}
	}
	if inputData.MachinesRecipesList != nil {
		for _, entry := range inputData.MachinesRecipesList {
			entry.UsersId = uint(userId)
		}
		result, err := h.MachineRecipeRepo.InsertMachinesRecipes(r.Context(), inputData.MachinesRecipesList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not insert requested machines_recipes data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.MachinesRecipesInserted = uint(noRows)
		}
	}
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("data has been inserted, but could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

func (h *CRUD) Update(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//jwt = token with dispatcher server secret key, id of user who received the token and issue date of the token, not optional
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
	inputData := JSONData{}
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("could not parse received body, reason: %w", err).Error()))
		return
	}
	var response struct {
		MachinesUpdated        uint
		ResourcesUpdated       uint
		RecipesUpdated         uint
		RecipesInputsUpdated   uint
		RecipesOutputsUpdated  uint
		MachinesRecipesUpdated uint
	}
	response.MachinesUpdated = 0
	response.ResourcesUpdated = 0
	response.RecipesUpdated = 0
	response.RecipesInputsUpdated = 0
	response.RecipesOutputsUpdated = 0
	response.MachinesRecipesUpdated = 0
	skipRows := false
	if inputData.MachinesList != nil {
		for _, entry := range inputData.MachinesList {
			entry.UsersId = uint(userId)
		}
		result, err := h.MachineRepo.UpdateMachines(r.Context(), inputData.MachinesList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not update requested machines data, reason: %w", err).Error()))
		}
		if !skipRows {
			for _, row := range result {
				noRows, err := row.RowsAffected()
				if err != nil {
					w.Write([]byte("database driver does not support returning numbers of rows affected"))
					skipRows = true
				}
				response.MachinesUpdated += uint(noRows)
			}
		}
	}
	if inputData.ResourcesList != nil {
		for _, entry := range inputData.ResourcesList {
			entry.UsersId = uint(userId)
		}
		result, err := h.ResourceRepo.UpdateResources(r.Context(), inputData.ResourcesList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not update requested resources data, reason: %w", err).Error()))
		}
		if !skipRows {
			for _, row := range result {
				noRows, err := row.RowsAffected()
				if err != nil {
					w.Write([]byte("database driver does not support returning numbers of rows affected"))
					skipRows = true
				}
				response.ResourcesUpdated += uint(noRows)
			}
		}
	}
	if inputData.RecipesList != nil {
		for _, entry := range inputData.RecipesList {
			entry.UsersId = uint(userId)
		}
		result, err := h.RecipeRepo.UpdateRecipes(r.Context(), inputData.RecipesList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not update requested recipes data, reason: %w", err).Error()))
		}
		if !skipRows {
			for _, row := range result {
				noRows, err := row.RowsAffected()
				if err != nil {
					w.Write([]byte("database driver does not support returning numbers of rows affected"))
					skipRows = true
				}
				response.RecipesUpdated += uint(noRows)
			}
		}
	}
	if inputData.RecipesInputsList != nil {
		for _, entry := range inputData.RecipesInputsList {
			entry.UsersId = uint(userId)
		}
		result, err := h.RecipeinputRepo.UpdateRecipesInputs(r.Context(), inputData.RecipesInputsList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not update requested recipes_inputs data, reason: %w", err).Error()))
		}
		if !skipRows {
			for _, row := range result {
				noRows, err := row.RowsAffected()
				if err != nil {
					w.Write([]byte("database driver does not support returning numbers of rows affected"))
					skipRows = true
				}
				response.RecipesInputsUpdated += uint(noRows)
			}
		}
	}
	if inputData.RecipesOutputsList != nil {
		for _, entry := range inputData.RecipesOutputsList {
			entry.UsersId = uint(userId)
		}
		result, err := h.RecipeoutputRepo.UpdateRecipesOutputs(r.Context(), inputData.RecipesOutputsList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not update requested recipes_outputs data, reason: %w", err).Error()))
		}
		if !skipRows {
			for _, row := range result {
				noRows, err := row.RowsAffected()
				if err != nil {
					w.Write([]byte("database driver does not support returning numbers of rows affected"))
					skipRows = true
				}
				response.RecipesOutputsUpdated += uint(noRows)
			}
		}
	}
	if inputData.MachinesRecipesList != nil {
		for _, entry := range inputData.MachinesRecipesList {
			entry.UsersId = uint(userId)
		}
		result, err := h.MachineRecipeRepo.UpdateMachinesRecipes(r.Context(), inputData.MachinesRecipesList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not update requested machines_recipes data, reason: %w", err).Error()))
		}
		if !skipRows {
			for _, row := range result {
				noRows, err := row.RowsAffected()
				if err != nil {
					w.Write([]byte("database driver does not support returning numbers of rows affected"))
					skipRows = true
				}
				response.MachinesRecipesUpdated += uint(noRows)
			}
		}
	}
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("data has been updated, but could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

func (h *CRUD) Delete(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//jwt = token with dispatcher server secret key, id of user who received the token and issue date of the token, not optional
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
	var inputData struct {
		MachinesIds        []int
		ResourcesIds       []int
		RecipesIds         []int
		RecipesInputsIds   []int
		RecipesOutputsIds  []int
		MachinesRecipesIds []int
	}
	var response struct {
		MachinesDeleted        uint
		ResourcesDeleted       uint
		RecipesDeleted         uint
		RecipesInputsDeleted   uint
		RecipesOutputsDeleted  uint
		MachinesRecipesDeleted uint
	}
	response.MachinesDeleted = 0
	response.ResourcesDeleted = 0
	response.RecipesDeleted = 0
	response.RecipesInputsDeleted = 0
	response.RecipesOutputsDeleted = 0
	response.MachinesRecipesDeleted = 0
	skipRows := false
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("could not parse received body, reason: %w", err).Error()))
		return
	}
	if inputData.MachinesIds != nil {
		result, err := h.MachineRepo.DeleteMachines(r.Context(), inputData.MachinesIds, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not delete requested machines data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.MachinesDeleted = uint(noRows)
		}
	}
	if inputData.ResourcesIds != nil {
		result, err := h.ResourceRepo.DeleteResources(r.Context(), inputData.ResourcesIds, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not delete requested resources data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.ResourcesDeleted = uint(noRows)
		}
	}
	if inputData.RecipesIds != nil {
		result, err := h.RecipeRepo.DeleteRecipes(r.Context(), inputData.RecipesIds, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not delete requested recipes data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.RecipesDeleted = uint(noRows)
		}
	}
	if inputData.RecipesInputsIds != nil {
		result, err := h.RecipeinputRepo.DeleteRecipesInputs(r.Context(), inputData.RecipesInputsIds, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not delete requested recipes_inputs data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.RecipesInputsDeleted = uint(noRows)
		}
	}
	if inputData.RecipesOutputsIds != nil {
		result, err := h.RecipeoutputRepo.DeleteRecipesOutputs(r.Context(), inputData.RecipesOutputsIds, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not delete requested recipes_outputs data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.RecipesOutputsDeleted = uint(noRows)
		}
	}
	if inputData.MachinesRecipesIds != nil {
		result, err := h.MachineRecipeRepo.DeleteMachinesRecipes(r.Context(), inputData.MachinesRecipesIds, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("could not delete requested recipes_outputs data, reason: %w", err).Error()))
		}
		if !skipRows {
			noRows, err := result.RowsAffected()
			if err != nil {
				w.Write([]byte("database driver does not support returning numbers of rows affected"))
				skipRows = true
			}
			response.MachinesRecipesDeleted = uint(noRows)
		}
	}
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("data has been deleted, but could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
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
