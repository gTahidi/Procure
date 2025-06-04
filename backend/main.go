package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"procurement/database"
	"procurement/handlers"
	"procurement/models"
	appMiddleware "procurement/middleware"
)

func main() {
	// Initialize Database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connection and initialization successful.")

	// Get DB instance for migration
	db := database.GetDB()

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Requisition{},
		&models.RequisitionItem{},
		&models.Tender{},
		&models.Bid{},
		&models.BidItem{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration successful.")

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

			// Register Tender routes
			tenderHandler := handlers.NewTenderHandler(db)
			authRouter.Post("/tenders", tenderHandler.CreateTender)
			authRouter.Get("/tenders", tenderHandler.GetTenders)
			authRouter.Get("/tenders/{id}", tenderHandler.GetTenderByID)
			authRouter.Put("/tenders/{id}", tenderHandler.UpdateTender)
			// authRouter.Delete("/tenders/{id}", tenderHandler.DeleteTender) // Commented out for now as DeleteTender is not yet implemented

			// Register Bid routes (New)
			bidHandler := handlers.NewBidHandler(db) // Create BidHandler instance
			// POST /api/tenders/{tender_id}/bids - Supplier creates a bid for a tender
			authRouter.Post("/tenders/{tenderId}/bids", bidHandler.CreateBid) // Placeholder for actual handler method

			// GET /api/tenders/{tender_id}/bids - Procurement officer lists bids for a tender
			authRouter.Get("/tenders/{tenderId}/bids", bidHandler.ListTenderBids) // Placeholder

			// GET /api/my-bids - Supplier lists their own submitted bids
			authRouter.Get("/my-bids", bidHandler.ListMyBids) // Placeholder

			// Add other authenticated routes here
		})

		// Public routes (if any)
		// r.Get("/some-public-route", somePublicHandler)
	})

	// Simple health check route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Procurement Backend!"))
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
