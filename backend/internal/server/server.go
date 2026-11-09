// Package server implements the calendar HTTP API.
package server

import (
	"encoding/json"
	"net/http"

	"github.com/godovdm/ai-for-developers-project-386/backend/internal/api/generated"
)

type server struct{}

// New returns the calendar API handler.
func New() http.Handler {
	return generated.Handler(server{})
}

func (server) GetHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(generated.HealthStatus{Status: "ok"})
}
