package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"procurement/database"
	"procurement/handlers"
	appMiddleware "procurement/middleware"
)

func main() {
	if err := database.InitDB(); err != nil { // Initialize DB connection pool
		log.Fatalf("FATAL: Could not initialize database: %v", err)
	}
	database.SetupDatabaseSchema() // Ensure schema is set up on startup using GORM

	r := chi.NewRouter()

	// CORS Middleware (assuming github.com/rs/cors)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"}, // Add your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"}, // Add any other headers your frontend sends
		AllowCredentials: true,
		Debug:            true, // Enable for debugging CORS issues
	})
	r.Use(c.Handler) // Apply the CORS middleware

	// Standard Chi middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // Log server requests
	r.Use(middleware.Recoverer)

	// API routes group
	r.Route("/api", func(apiRouter chi.Router) {
		// Register user routes (e.g., /api/users/sync) - these might be public or have their own auth
		handlers.RegisterUserRoutes(apiRouter)

		// Authenticated routes group
		apiRouter.Group(func(authRouter chi.Router) {
			authRouter.Use(appMiddleware.TokenMiddleware) // Apply the token middleware

			// Register requisition routes
			authRouter.Post("/requisitions", handlers.CreateRequisitionHandler)
			authRouter.Get("/requisitions", handlers.ListRequisitionsHandler)
			authRouter.Get("/requisitions/{id}", handlers.GetRequisitionHandler) // New route for single requisition
			// Add other authenticated requisition routes here (GET, PUT, DELETE)
		})
	})

	// Simple health check route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Procurement Backend!"))
	})

	port := ":8080"
	log.Printf("Backend server starting on port %s...", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
