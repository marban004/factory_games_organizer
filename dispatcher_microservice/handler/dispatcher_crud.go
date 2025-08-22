package handler

import "net/http"

type DispatcherCrud struct {
	CommonHandlerFunctions     CommonHandlerFunctions
	CrudMicroservicesAddresses []string
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
//	@Success		200	{object}	handler.JSONDataCrud
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/crud/selectbyid [get]
//
//	@Security		apiTokenAuth
func (h *DispatcherCrud) SelectByID(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, _ := h.CommonHandlerFunctions.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	h.CommonHandlerFunctions.redirectRequest(w, r, "selectbyid", h.CrudMicroservicesAddresses)
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
//	@Success		200	{object}	handler.JSONDataCrud
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/crud/select [get]
//
//	@Security		apiTokenAuth
func (h *DispatcherCrud) Select(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, _ := h.CommonHandlerFunctions.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	h.CommonHandlerFunctions.redirectRequest(w, r, "select", h.CrudMicroservicesAddresses)
}

// Insert insert record(s) into the database
//
//	@Description	Insert data into database. The user to whom the ownership of records is assigned is the user who presented the authentication token.
//	@Param			insert	body	handler.JSONDataCrud	true	"Data to be inserted into database"
//	@Tags			CRUD Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.InsertResponseCrud
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/crud [post]
//
//	@Security		apiTokenAuth
func (h *DispatcherCrud) Insert(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, _ := h.CommonHandlerFunctions.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	h.CommonHandlerFunctions.redirectRequest(w, r, "", h.CrudMicroservicesAddresses)
}

// Update update record(s) in the database
//
//	@Description	Updates data in database. Updates the records based on "id" field of an element in the array sent in request body. If a record with a particular id does not belong to the user who presented authentication token, then that record is not updated.
//	@Param			update	body	handler.JSONDataCrud	true	"Data to be updated in the database"
//	@Tags			CRUD Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.UpdateResponseCrud
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/crud [put]
//
//	@Security		apiTokenAuth
func (h *DispatcherCrud) Update(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, _ := h.CommonHandlerFunctions.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	h.CommonHandlerFunctions.redirectRequest(w, r, "", h.CrudMicroservicesAddresses)
}

// Delete delete record(s) in the database
//
//	@Description	Deletes data in database. Each table has it's own id list to be deleted. If a record with a particular id does not belong to the user who presented authentication token, then that record is not deleted.
//	@Param			delete	body	handler.DeleteInputCrud	true	"Data to be deleted in the database"
//	@Tags			CRUD Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.DeleteResponseCrud
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/crud [delete]
//
//	@Security		apiTokenAuth
func (h *DispatcherCrud) Delete(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, _ := h.CommonHandlerFunctions.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	h.CommonHandlerFunctions.redirectRequest(w, r, "", h.CrudMicroservicesAddresses)
}

// Delete delete record(s) in the database
//
//	@Description	Deletes all data in the database that belongs to user who presented the authentication token.
//	@Tags			CRUD Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.DeleteResponseCrud
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/crud/user [delete]
//
//	@Security		apiTokenAuth
func (h *DispatcherCrud) DeleteByUser(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, _ := h.CommonHandlerFunctions.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	h.CommonHandlerFunctions.redirectRequest(w, r, "user", h.CrudMicroservicesAddresses)
}
