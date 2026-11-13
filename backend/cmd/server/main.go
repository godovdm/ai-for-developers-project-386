// Package main starts the calendar HTTP server.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/godovdm/ai-for-developers-project-386/backend/internal/server"
	"github.com/godovdm/ai-for-developers-project-386/backend/internal/web"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/", server.New())
	dir := os.Getenv("FRONTEND_DIR")
	if dir == "" {
		dir = "../frontend/dist"
	}
	mux.Handle("/", web.SPA(dir))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Print("starting server")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
