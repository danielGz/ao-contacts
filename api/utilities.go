package api

import (
	"encoding/json"
	"net/http"
)

func JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Type for all responses from this handler
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func jsonResponseError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
