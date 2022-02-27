package handlers

import (
	"encoding/json"
	"net/http"
)

type v1Error struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func errorMsg(w http.ResponseWriter, message string, statusCode int) {
	err := v1Error{
		Status:  statusCode,
		Message: message,
	}
	errorJson, _ := json.Marshal(err)
	http.Error(w, string(errorJson), statusCode)
}

func invalidMethodError(w http.ResponseWriter) {
	errorMsg(w, "invalid method", http.StatusMethodNotAllowed)
}