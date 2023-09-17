package handler

import (
	"net/http"
	"time"
)

// handleHealth is a basic health route handler.
func (*handler) handleHealth() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Disable caching
		w.Header().Set("Expires", time.Unix(0, 0).Format(time.RFC1123))
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("X-Accel-Expires", "0")

		w.WriteHeader(http.StatusOK)
	}
}
