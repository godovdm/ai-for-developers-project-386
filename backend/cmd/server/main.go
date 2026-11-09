// Package main starts the calendar HTTP server.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/godovdm/ai-for-developers-project-386/backend/internal/server"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/", server.New())
	dir := os.Getenv("FRONTEND_DIR")
	if dir == "" {
		dir = "../frontend/dist"
	}
	mux.Handle("/", http.FileServer(http.Dir(dir)))
	log.Print("listening on http://localhost:8080")
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
