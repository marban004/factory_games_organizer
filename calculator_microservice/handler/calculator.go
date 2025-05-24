package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	microservicelogiccalculator "github.com/marban004/factory_games_organizer/calculator_microservice/microservice_logic_calculator"
)

type Calculator struct {
	DB *sql.DB
}

func (c *Calculator) Calculate(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//userid = user with whose data we want to generate production tree, not optional
	//resource = resource for which we want to generate production tree, not optional
	//rate = target production rate per second for specified resource, not optional
	//alt_recipes = recipes that can be used besides default recipes, optional
	//alt_machines = machines that can be used besides default machines, optional
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
		w.Write([]byte("userid should be a positive floating point number and cannot be empty"))
		return
	}
	recipes_names := []string{r.URL.Query().Get("alt_recipes")}
	machine_names := []string{r.URL.Query().Get("alt_machines")}
	byteJSONRepresentation, err := microservicelogiccalculator.Calculate(r.Context(), userId, desiredResourceName, float32(desiredRate), recipes_names, machine_names, c.DB)
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
