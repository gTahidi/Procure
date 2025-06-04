package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// SanitizeFilename removes or replaces characters that are problematic for filenames.
// This version replaces known problematic characters with underscores.
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9_.-]`) // Allow letters, numbers, underscore, dot, hyphen

func SanitizeFilename(filename string) string {
	// Replace spaces with underscores
	cleaned := strings.ReplaceAll(filename, " ", "_")
	// Replace any characters not in the allowed set (alphanumeric, _, ., -) with an underscore
	cleaned = nonAlphanumericRegex.ReplaceAllString(cleaned, "_")
	// Reduce multiple underscores to a single underscore
	cleaned = regexp.MustCompile(`_+`).ReplaceAllString(cleaned, "_")
	// Remove leading/trailing underscores that might result from replacements
	cleaned = strings.Trim(cleaned, "_")

	// Prevent empty filenames or filenames that are just "." or ".."
	if cleaned == "" || cleaned == "." || cleaned == ".." {
		return "sanitized_filename"
	}
	return cleaned
}

// getFormValuePointer retrieves a form value by key.
// If the key is not present or the value is empty, it returns nil.
// Otherwise, it returns a pointer to the string value.
func getFormValuePointer(r *http.Request, key string) *string {
	val := r.FormValue(key)
	if val == "" {
		return nil
	}
	return &val
}

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
