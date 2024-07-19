package utils

import (
	"encoding/json"
	"gitbeam/models"
	"net/http"
)

func WriteHTTPError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(&models.Result{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	})
}

func WriteHTTPSuccess(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&models.Result{
		Success: true,
		Message: message,
		Data:    data,
	})
}
