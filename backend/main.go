package main

import (
	"log"
	"net/http"
	"procurement/database"
	"procurement/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("Initializing database...")
	if err := database.InitDB(); err != nil {
		log.Fatalf("FATAL: Could not initialize database: %v", err)
	}
	log.Println("Database initialized successfully.")

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // Logs the request path, method, duration, etc.
	r.Use(middleware.Recoverer) // Recovers from panics and returns a 500 error

	// Basic health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Procurement Backend!"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/requisitions", handlers.CreateRequisitionHandler)
		// Add other requisition routes here (GET, PUT, DELETE)
	})

	port := ":8080"
	log.Printf("Backend server starting on port %s...", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("FATAL: Could not start server on port %s: %v", port, err)
	}
}
