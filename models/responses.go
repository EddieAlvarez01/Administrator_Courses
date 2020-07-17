package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//NewResponseJSON GENERIC STRUCT TO RESPOND TO THE SERVER
func NewResponseJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	response := response{Code: code, Message: message, Data: data}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in encoding response struct"), http.StatusInternalServerError)
		return
	}
}