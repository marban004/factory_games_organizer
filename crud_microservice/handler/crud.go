package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	custommiddleware "github.com/marban004/factory_games_organizer/custom_middleware"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/model"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine"
	machinerecipe "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine_recipe"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe"
	recipeinput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_input"
	recipeoutput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_output"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/resource"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type JSONData struct {
	MachinesList        []model.MachineInfo
	ResourcesList       []model.ResourceInfo
	RecipesList         []model.RecipeInfo
	RecipesInputsList   []model.RecipeInputOutputInfo
	RecipesOutputsList  []model.RecipeInputOutputInfo
	MachinesRecipesList []model.MachinesRecipesInfo
}

type InsertResponse struct {
	MachinesInserted        uint
	ResourcesInserted       uint
	RecipesInserted         uint
	RecipesInputsInserted   uint
	RecipesOutputsInserted  uint
	MachinesRecipesInserted uint
}

type UpdateResponse struct {
	MachinesUpdated        uint
	ResourcesUpdated       uint
	RecipesUpdated         uint
	RecipesInputsUpdated   uint
	RecipesOutputsUpdated  uint
	MachinesRecipesUpdated uint
}

type DeleteInput struct {
	MachinesIds        []int
	ResourcesIds       []int
	RecipesIds         []int
	RecipesInputsIds   []int
	RecipesOutputsIds  []int
	MachinesRecipesIds []int
}

type DeleteResponse struct {
	MachinesDeleted        uint
	ResourcesDeleted       uint
	RecipesDeleted         uint
	RecipesInputsDeleted   uint
	RecipesOutputsDeleted  uint
	MachinesRecipesDeleted uint
}

type CRUD struct {
	MachineRepo       *machine.MySQLRepo
	ResourceRepo      *resource.MySQLRepo
	RecipeRepo        *recipe.MySQLRepo
	RecipeinputRepo   *recipeinput.MySQLRepo
	RecipeoutputRepo  *recipeoutput.MySQLRepo
	MachineRecipeRepo *machinerecipe.MySQLRepo
	Secret            []byte
	StatTracker       *custommiddleware.DefaultApiStatTracker
}

type HealthResponse struct {
	MicroserviceStatus string
	DatabaseStatus     string
}

type StatsResponse struct {
	ApiUsageStats    *orderedmap.OrderedMap[string, map[string]int]
	TrackingPeriodMs int64
	NoPeriods        uint64
}

// Health return the status of microservice and associated database
//
//	@Description	Return the status of microservice and it's database. Default working state is signified by status "up".
//	@Tags			CRUD
//	@Success		200	{object}	handler.HealthResponse
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/health [get]
func (h *CRUD) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		MicroserviceStatus: "up",
	}
	err := h.MachineRepo.DB.PingContext(r.Context())
	if err != nil {
		response.DatabaseStatus = "connection disrupted"
	} else {
		response.DatabaseStatus = "up"
	}
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

// Stats return the usage stats of microservice
//
//	@Description	Return the usage stats of microservice.
//	@Tags			CRUD
//	@Success		200	{object}	handler.StatsResponse
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/stats [get]
func (h *CRUD) Stats(w http.ResponseWriter, r *http.Request) {
	endpointResponse := StatsResponse{ApiUsageStats: h.StatTracker.GetStats(), TrackingPeriodMs: h.StatTracker.Period, NoPeriods: h.StatTracker.MaxLen}
	byteJSONRepresentation, err := json.Marshal(endpointResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

// SelectByID return the record(s) from database
//
//	@Description	Return the records from database specified by id. Id(s) is specified for each table in the database. If an id parameter for a particular table is omitted, the records are not retreived from that table. Each parameter can be present multiple times, in which case all records from a particular table, with those ids will be retreived and returned in an array. Data is returned for the user that provided authentication token.
//	@Param			machines_id			query	integer	false	"Id of machines to be retreived from database"
//	@Param			resources_id		query	integer	false	"Id of resources to be retreived from database"
//	@Param			recipes_id			query	integer	false	"Id of recipes to be retreived from database"
//	@Param			recipes_inputs_id	query	integer	false	"Id of recipes inputs to be retreived from database"
//	@Param			recipes_outputs_id	query	integer	false	"Id of recipes outputs to be retreived from database"
//	@Param			machines_recipes_id	query	integer	false	"Id of machines recipes to be retreived from database"
//	@Tags			CRUD Authorization required
//	@Success		200	{object}	handler.JSONData
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/selectbyid [get]
//
//	@Security		apiTokenAuth
func (h *CRUD) SelectByID(w http.ResponseWriter, r *http.Request) {
	//parameters that are not mentioned in swagger directly:
	//jwt = token with dispatcher server secret key, id of user who received the token and issue date of the token, not optional
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, userId := h.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
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

// Select return the record(s) from database
//
//	@Description	Return the records from database specified by id range. Start of range and it's size is specified for each table separately. Size describes number of records to be returned. Ranges include the starting id. Records are only returned for the user that presented authentication token. If start of the range is missing for particular table, then it is assumed to be 1. If size is ommitted, then all records are retreived. If the start of a range is a record belonging to another user, then next record belonging to the user that presented a token is retreived instead.
//	@Param			machines_id_start			query	integer	false	"Id of first record to be retreived from machines table"
//	@Param			machines_rows				query	integer	false	"Number of rows to be returned from machines table"
//	@Param			resources_id_start			query	integer	false	"Id of first record to be retreived from resources table"
//	@Param			resources_rows				query	integer	false	"Number of rows to be returned from resources table"
//	@Param			recipes_id_start			query	integer	false	"Id of first record to be retreived from recipes table"
//	@Param			recipes_rows				query	integer	false	"Number of rows to be returned from recipes table"
//	@Param			recipes_inputs_id_start		query	integer	false	"Id of first record to be retreived from recipes_inputs table"
//	@Param			recipes_inputs_rows			query	integer	false	"Number of rows to be returned from recipes_inputs table"
//	@Param			recipes_outputs_id_start	query	integer	false	"Id of first record to be retreived from recipes_outputs table"
//	@Param			recipes_outputs_rows		query	integer	false	"Number of rows to be returned from recipes_outputs table"
//	@Param			machines_recipes_id_start	query	integer	false	"Id of first record to be retreived from machines_recipes table"
//	@Param			machines_recipes_rows		query	integer	false	"Number of rows to be returned from machines_recipes table"
//	@Tags			CRUD Authorization required
//	@Success		200	{object}	handler.JSONData
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/select [get]
//
//	@Security		apiTokenAuth
func (h *CRUD) Select(w http.ResponseWriter, r *http.Request) {
	//parameters that are not mentioned in swagger directly:
	//jwt = token with dispatcher server secret key, id of user who received the token and issue date of the token, not optional
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, userId := h.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
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

// Insert insert record(s) into the database
//
//	@Description	Insert data into database. The user to whom the ownership of records is assigned is the user who presented the authentication token.
//	@Param			insert	body	handler.JSONData	true	"Data to be inserted into database"
//	@Tags			CRUD Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.InsertResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/ [post]
//
//	@Security		apiTokenAuth
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
		w.WriteHeader(http.StatusUnauthorized)
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
	response := InsertResponse{}
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
	w.WriteHeader(http.StatusCreated)
	w.Write(byteJSONRepresentation)
}

// Update update record(s) in the database
//
//	@Description	Updates data in database. Updates the records based on "id" field of an element in the array sent in request body. If a record with a particular id does not belong to the user who presented authentication token, then that record is not updated.
//	@Param			update	body	handler.JSONData	true	"Data to be updated in the database"
//	@Tags			CRUD Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.UpdateResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/ [put]
//
//	@Security		apiTokenAuth
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
		w.WriteHeader(http.StatusUnauthorized)
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
	response := UpdateResponse{}
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
			return
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
			return
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
			return
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
			return
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
			return
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
			return
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

// Delete delete record(s) in the database
//
//	@Description	Deletes data in database. Each table has it's own id list to be deleted. If a record with a particular id does not belong to the user who presented authentication token, then that record is not deleted.
//	@Param			delete	body	handler.DeleteInput	true	"Data to be deleted in the database"
//	@Tags			CRUD Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.DeleteResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/ [delete]
//
//	@Security		apiTokenAuth
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
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	inputData := DeleteInput{}
	response := DeleteResponse{}
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

// Delete delete record(s) in the database
//
//	@Description	Deletes all data in the database that belongs to user who presented the authentication token.
//	@Tags			CRUD Authorization required
//
//	@Success		200	{object}	handler.DeleteResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/user [delete]
//
//	@Security		apiTokenAuth
func (h *CRUD) DeleteByUser(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	response := DeleteResponse{}
	response.MachinesDeleted = 0
	response.ResourcesDeleted = 0
	response.RecipesDeleted = 0
	response.RecipesInputsDeleted = 0
	response.RecipesOutputsDeleted = 0
	response.MachinesRecipesDeleted = 0
	ctx := r.Context()
	transaction, err := h.MachineRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not start a transaction, reason: %w", err).Error()))
		return
	}
	skipRows := false

	result, err := h.MachineRepo.DeleteMachinesByUserId(ctx, transaction, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not delete requested machines data, reason: %w", err).Error()))
		return
	}
	if !skipRows {
		noRows, err := result.RowsAffected()
		if err != nil {
			w.Write([]byte("database driver does not support returning numbers of rows affected"))
			skipRows = true
		}
		response.MachinesDeleted = uint(noRows)
	}
	result, err = h.ResourceRepo.DeleteResourcesByUserId(ctx, transaction, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not delete requested resources data, reason: %w", err).Error()))
		return
	}
	if !skipRows {
		noRows, err := result.RowsAffected()
		if err != nil {
			w.Write([]byte("database driver does not support returning numbers of rows affected"))
			skipRows = true
		}
		response.ResourcesDeleted = uint(noRows)
	}
	result, err = h.RecipeRepo.DeleteRecipesByUserId(ctx, transaction, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not delete requested recipes data, reason: %w", err).Error()))
		return
	}
	if !skipRows {
		noRows, err := result.RowsAffected()
		if err != nil {
			w.Write([]byte("database driver does not support returning numbers of rows affected"))
			skipRows = true
		}
		response.RecipesDeleted = uint(noRows)
	}
	result, err = h.RecipeinputRepo.DeleteRecipesInputsByUserId(ctx, transaction, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not delete requested recipes_inputs data, reason: %w", err).Error()))
		return
	}
	if !skipRows {
		noRows, err := result.RowsAffected()
		if err != nil {
			w.Write([]byte("database driver does not support returning numbers of rows affected"))
			skipRows = true
		}
		response.RecipesInputsDeleted = uint(noRows)
	}
	result, err = h.RecipeoutputRepo.DeleteRecipesOutputsByUserId(ctx, transaction, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not delete requested recipes_outputs data, reason: %w", err).Error()))
		return
	}
	if !skipRows {
		noRows, err := result.RowsAffected()
		if err != nil {
			w.Write([]byte("database driver does not support returning numbers of rows affected"))
			skipRows = true
		}
		response.RecipesOutputsDeleted = uint(noRows)
	}
	result, err = h.MachineRecipeRepo.DeleteMachinesRecipesByUserId(ctx, transaction, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not delete requested machines_inputs data, reason: %w", err).Error()))
		return
	}
	if !skipRows {
		noRows, err := result.RowsAffected()
		if err != nil {
			w.Write([]byte("database driver does not support returning numbers of rows affected"))
			skipRows = true
		}
		response.MachinesRecipesDeleted = uint(noRows)
	}
	err = transaction.Commit()
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("an error occurred, could not rollback transaction, reason: %w", rollbackErr).Error()))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Errorf("an error occurred, transaction has been rolled back, reason: %w", err).Error()))
			return
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
func (h *CRUD) verifyJWT(jwtString string) (bool, int) {
	token, err := jwt.Parse(jwtString, func(*jwt.Token) (interface{}, error) {
		return h.Secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return false, 0
	}
	if !token.Valid {
		return false, 0
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"].(float64)
	expTime := int64(claims["exp"].(float64))
	issueTime := int64(claims["iat"].(float64))
	if time.Now().Unix() > expTime {
		return false, 0
	}
	if time.Now().Unix() <= issueTime {
		return false, 0
	}
	return true, int(userId)
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
