//     This is Factory Games Organizer api. Api is responsible for creating, updating and authenicating api users, CRUD operations on database associated with the api and provides production calculator service.
//     Copyright (C) 2025  Marek Banaś

//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.

//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see https://www.gnu.org/licenses/.

package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	custommiddleware "github.com/marban004/factory_games_organizer/custom_middleware"
	microservicelogiccalculator "github.com/marban004/factory_games_organizer/microservice_logic_calculator"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Calculator struct {
	DB          *sql.DB
	StatTracker *custommiddleware.DefaultApiStatTracker
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

// Calculate return the calculated production tree for specified resource
//
//	@Description	Calculate the machines and resources needed to produce target resource with provided production rate per second. Alternative Recipe and Alternative Machine parameters can be present multiple times in request query.
//	@Param			userid		query	string	true	"Id of users whose data will be used as the base for calculation"
//	@Param			resource	query	string	true	"Resource to be produced"
//	@Param			rate		query	string	true	"Target production rate for the specified resource"
//	@Param			alt_recipe	query	string	false	"Alternative recipe to take into consideration when calculating production tree"
//	@Param			alt_machine	query	string	false	"Alternative machine to take into consideration when calculating production tree"
//	@Tags			Calculator
//	@Success		200	{object}	microservicelogiccalculator.ProductionTree
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/calculate [get]
func (h *Calculator) Calculate(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("userid"))
	if err != nil || userId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("userid should be a positive integer and cannot be empty"))
		return
	}
	desiredResourceName := r.URL.Query().Get("resource")
	if desiredResourceName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("resource parameter cannot be empty"))
		return
	}
	desiredRate, err := strconv.ParseFloat(r.URL.Query().Get("rate"), 32)
	if err != nil || desiredRate <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("rate should be a positive floating point number and cannot be empty"))
		return
	}
	recipes_names := r.URL.Query()["alt_recipe"]
	machine_names := r.URL.Query()["alt_machine"]
	byteJSONRepresentation, err := microservicelogiccalculator.Calculate(r.Context(), userId, desiredResourceName, float32(desiredRate), recipes_names, machine_names, h.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not generate production tree for '%s', reason: %w", desiredResourceName, err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
	// test url 192.168.31.74:3000/calculate?userid=1&resource=reinforced_iron_plate&rate=0.5
	// w.Write([]byte("works maybe"))
}

// Health return the status of microservice and associated database
//
//	@Description	Return the status of microservice and it's database. Default working state is signified by status "up".
//	@Tags			Calculator
//	@Success		200	{object}	handler.HealthResponse
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/health [get]
func (h *Calculator) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		MicroserviceStatus: "up",
	}
	err := h.DB.PingContext(r.Context())
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
//	@Tags			Calculator
//	@Success		200	{object}	handler.StatsResponse
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/stats [get]
func (h *Calculator) Stats(w http.ResponseWriter, r *http.Request) {
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
