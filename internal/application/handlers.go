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
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func (a *Application) calculateHandler(w http.ResponseWriter, r *http.Request) {
	wrapped := w.(*responseWriter)

	if r.Method != http.MethodPost {
		wrapped.error = "метод не разрешен"
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wrapped.error = "некорректное тело запроса"
		http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	result, err := calculation.Calc(req.Expression)
	response := CalculationResponse{}

	if err != nil {
		wrapped.error = err.Error()
		response.Error = "Expression is not valid"
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		response.Result = result
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		wrapped.error = "ошибка сериализации"
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
}
