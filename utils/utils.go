package utils

import (
	"bytes"
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

func UnPack(in interface{}, target interface{}) error {
	var e1 error
	var b []byte
	switch in := in.(type) {
	case []byte:
		b = in
	// Do something.
	default:
		// Do the rest.
		b, e1 = json.Marshal(in)
		if e1 != nil {
			return e1
		}
	}

	buf := bytes.NewBuffer(b)
	enc := json.NewDecoder(buf)
	enc.UseNumber()
	if err := enc.Decode(&target); err != nil {
		return err
	}
	return nil
}
