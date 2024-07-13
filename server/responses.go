package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
	"net/http"
	"time"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError))) //nolint:all
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes) //nolint:all
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes) //nolint:all
}

// HealthResponse represents the response for the health check.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type GameResultResponse struct {
	ID                int                      `json:"id"`
	UserID            uuid.UUID                `json:"userId"`
	GameStatus        entity.GameStatus        `json:"state"`
	TransactionSource entity.TransactionSource `json:"source"`
	TransactionID     string                   `json:"transactionId"`
	Amount            float64                  `json:"amount"`
	CreatedAt         time.Time                `json:"createdAt"`
}
