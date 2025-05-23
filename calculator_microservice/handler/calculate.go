package handler

import (
	"net/http"
)

type Calculator struct{}

func (c *Calculator) Calculate(w http.ResponseWriter, r *http.Request) {
	// byteJSONRepresentation, err := microservicelogiccalculator.Calculate()
	// if err != nil {
	// 	fmt.Printf("Could not generate production tree for '%s', reason: %v \n", desiredResourceName, err)
	// }
	// w.Write(byteJSONRepresentation)
	w.Write([]byte("works maybe"))
}
