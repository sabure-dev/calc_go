package application

import (
	"encoding/json"
	"net/http"

	calculation "github.com/sabure-dev/calc_go/pkg/calculation"
)

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func (a *Application) calculateHandler(w http.ResponseWriter, r *http.Request) {
	wrapped := w.(*responseWriter)

	if r.Method != http.MethodPost {
		wrapped.error = "неверный метод запроса"
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wrapped.error = "некорректное тело запроса: " + err.Error()
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	wrapped.expression = req.Expression
	result, err := calculation.Calc(req.Expression)
	response := CalculationResponse{}

	if err != nil {
		wrapped.error = "ошибка вычисления: " + err.Error()
		response.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		response.Result = result
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		wrapped.error = "ошибка сериализации ответа: " + err.Error()
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
