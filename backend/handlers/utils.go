package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithError returns a JSON error response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON returns a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("ERROR: Failed to encode response: %v\n", err)
		// If headers are already sent, we can't change the status code here.
		// Logging the error is important.
	}
}
