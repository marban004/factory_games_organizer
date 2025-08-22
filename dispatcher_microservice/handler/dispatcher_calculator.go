package handler

import "net/http"

type DispatcherCalculator struct {
	CommonHandlerFunctions           CommonHandlerFunctions
	CalculatorMicroservicesAddresses []string
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
//	@Success		200	{object}	handler.ProductionTreeCalculator
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/calculator/calculate [get]
func (h *DispatcherCalculator) Calculate(w http.ResponseWriter, r *http.Request) {
	h.CommonHandlerFunctions.redirectRequest(w, r, "calculate", h.CalculatorMicroservicesAddresses)
}
