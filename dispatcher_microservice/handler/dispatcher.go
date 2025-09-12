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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	custommiddleware "github.com/marban004/factory_games_organizer/custom_middleware"
)

type Dispatcher struct {
	UsersMicroservicesAddresses      []string
	CrudMicroservicesAddresses       []string
	CalculatorMicroservicesAddresses []string
	StatTracker                      *custommiddleware.DefaultApiStatTracker
}

// Health return the status of microservices and their associated databases
//
//	@Description	Return the status of microservice and it's database. Default working state is signified by status "up".
//	@Tags			Dispatcher
//	@Success		200	{object}	handler.HealthResponse
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/health [get]
func (h *Dispatcher) Health(w http.ResponseWriter, r *http.Request) {
	endpointResponse := HealthResponse{
		DispatcherStatus:       "up",
		UsersMicroservice:      []MicroserviceHealth{},
		CrudMicroservice:       []MicroserviceHealth{},
		CalculatorMicroservice: []MicroserviceHealth{},
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	err := h.checkMicroservicesHealth(r.Context(), &endpointResponse.UsersMicroservice, client, h.UsersMicroservicesAddresses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not create request to microservice"))
		return
	}
	err = h.checkMicroservicesHealth(r.Context(), &endpointResponse.CrudMicroservice, client, h.CrudMicroservicesAddresses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not create request to microservice"))
		return
	}
	err = h.checkMicroservicesHealth(r.Context(), &endpointResponse.CalculatorMicroservice, client, h.CalculatorMicroservicesAddresses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not create request to microservice"))
		return
	}

	byteJSONRepresentation, err := json.Marshal(endpointResponse)
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
//	@Tags			Dispatcher
//	@Success		200	{object}	handler.StatsResponse
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/stats [get]
func (h *Dispatcher) Stats(w http.ResponseWriter, r *http.Request) {
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

func (h *Dispatcher) checkMicroservicesHealth(ctx context.Context, endpointResponseArray *[]MicroserviceHealth, client *http.Client, microserviceAddressList []string) error {
	for _, address := range microserviceAddressList {
		microserviceStatus := MicroserviceHealth{MicroserviceURL: address}
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/health", address), nil)
		if err != nil {
			return err
		}
		response, err := client.Do(request)
		if err != nil {
			microserviceStatus.MicroserviceStatus = "down"
			microserviceStatus.DatabaseStatus = "unknown"
		} else if response.StatusCode < 200 || response.StatusCode >= 300 {
			microserviceStatus.MicroserviceStatus = "health endpoint malfunction"
			microserviceStatus.DatabaseStatus = "unknown"
		} else {
			var microserviceResponse struct {
				MicroserviceStatus string
				DatabaseStatus     string
			}
			err := json.NewDecoder(response.Body).Decode(&microserviceResponse)
			if err != nil {
				microserviceStatus.MicroserviceStatus = "malformed response from microservice"
				microserviceStatus.DatabaseStatus = "unknown"
			} else {
				microserviceStatus.MicroserviceStatus = microserviceResponse.MicroserviceStatus
				microserviceStatus.DatabaseStatus = microserviceResponse.DatabaseStatus
			}
		}
		*endpointResponseArray = append(*endpointResponseArray, microserviceStatus)
	}
	return nil
}

// func (h *Users) convertArrToInt(input []string) []int {
// 	result := []int{}
// 	for _, value := range input {
// 		intValue, err := strconv.Atoi(value)
// 		if err != nil {
// 			continue
// 		}
// 		result = append(result, intValue)
// 	}
// 	return result
// }
