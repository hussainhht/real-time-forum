package utilities

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, data interface{}) error { //Decodes JSON from the request body.
	return json.NewDecoder(r.Body).Decode(data)
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error { //Encodes the provided data.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ErrorJSON(w http.ResponseWriter, status int, message string) { //Sends a JSON error response.
	WriteJSON(w, status, map[string]string{"error": message})
}