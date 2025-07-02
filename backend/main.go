package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"procurement/database"
	"procurement/handlers"
	appMiddleware "procurement/middleware"
	"procurement/models"
)

// serveFrontend serves the static SvelteKit application.
// It uses a catch-all route that first checks for a static file,
// and if not found, serves the index.html file for SPA routing.
func serveFrontend(r *chi.Mux, staticPath string) {
	fs := http.FileServer(http.Dir(staticPath))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Check if a file exists at the requested path
		filePath := filepath.Join(staticPath, r.URL.Path)
		stat, err := os.Stat(filePath)

		// If the file exists and it is NOT a directory, serve it.
		if err == nil && !stat.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		// For any other case (file not found, it's a directory, etc.),
		// serve the main index.html file to let the SPA router handle it.
		http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
	})
}

func main() {
	// --- DATABASE INITIALIZATION (No changes needed here) ---
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connection and initialization successful.")
	db := database.GetDB()
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

	// --- ROUTER & MIDDLEWARE SETUP (No changes needed here) ---
	r := chi.NewRouter()

	// Note on CORS: When serving from the same origin, CORS is not strictly necessary.
	// However, it's useful for local development (e.g., Vite dev server at :5173 hitting API at :8080).
	// For production, you could restrict this or remove it.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		Debug:            true,
	})
	r.Use(c.Handler)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- ROUTING LOGIC ---

	// **STEP 1: Mount all API routes first.**
	// All requests to /api/... will be handled by this group.
	r.Route("/api", func(apiRouter chi.Router) {
		handlers.RegisterUserRoutes(apiRouter)
		apiRouter.Group(func(authRouter chi.Router) {
			authRouter.Use(appMiddleware.TokenMiddleware)
			// ... all your existing authenticated routes go here ...
			authRouter.Post("/requisitions", handlers.CreateRequisitionHandler)
			authRouter.Get("/requisitions", handlers.ListRequisitionsHandler)
			authRouter.Get("/requisitions/{id}", handlers.GetRequisitionHandler)
			authRouter.Post("/requisitions/{id}/action", handlers.HandleRequisitionAction)

			tenderHandler := handlers.NewTenderHandler(db)
			authRouter.Post("/tenders", tenderHandler.CreateTender)
			// ... etc ...

			bidHandler := handlers.NewBidHandler(db)
			authRouter.Post("/tenders/{tenderId}/bids", bidHandler.CreateBid)
			// GET /api/tenders/{tender_id}/bids - Procurement officer lists bids for a tender
			authRouter.Get("/tenders/{tenderId}/bids", bidHandler.ListTenderBids) // Placeholder

			// GET /api/my-bids - Supplier lists their own submitted bids
			authRouter.Get("/my-bids", bidHandler.ListMyBids) // Placeholder

			// Register dashboard routes
			authRouter.Get("/dashboard/requisition-stats", handlers.GetRequisitionStatsHandler)
			authRouter.Get("/dashboard/recent-requisitions", handlers.GetRecentRequisitionsHandler)
			authRouter.Get("/dashboard/live-tenders", handlers.GetLiveTendersHandler)
			authRouter.Get("/dashboard/creation-rate", handlers.GetCreationRateHandler)

			// Routes for requester-specific dashboard data
			authRouter.Get("/dashboard/my-stats", handlers.GetMyRequisitionStatsHandler)
			authRouter.Get("/dashboard/my-recent-requisitions", handlers.GetMyRecentRequisitionsHandler)

			// Route for supplier-specific dashboard data
			authRouter.Get("/dashboard/supplier", handlers.GetSupplierDashboardDataHandler)

		})
	})

	// **STEP 2: Mount the frontend file server as a catch-all.**
	// This will handle any request that did not match /api.
	// The path "frontend/dist" matches the destination in your Dockerfile.
	serveFrontend(r, "frontend/dist")

	// --- START SERVER (No changes needed here) ---
	log.Println("Server starting on :8080, serving API and Frontend")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
