package http

import (
	"encoding/json"
	"net/http"
)

func NewErrorResponse(w http.ResponseWriter, statusCode int, message string) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(&ErrorResponse{Message: message})
}
