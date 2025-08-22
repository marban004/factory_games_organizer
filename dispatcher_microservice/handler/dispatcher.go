package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Dispatcher struct {
	UsersMicroservicesAddresses      []string
	CrudMicroservicesAddresses       []string
	CalculatorMicroservicesAddresses []string
}

// Health return the status of microservices and their associated databases
//
//	@Description	Return the status of microservice and it's database. Default working state is signified by status "up".
//	@Tags			All microservices
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

func (h *Dispatcher) checkMicroservicesHealth(ctx context.Context, endpointResponseArray *[]MicroserviceHealth, client *http.Client, microserviceAddressList []string) error {
	for _, address := range microserviceAddressList {
		microserviceStatus := MicroserviceHealth{MicroserviceURl: address}
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
