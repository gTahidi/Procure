package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"procurement/database"
	"procurement/handlers"
	appMiddleware "procurement/middleware"
	"procurement/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func serveFrontend(r *chi.Mux, staticPath string) {
	fs := http.FileServer(http.Dir(staticPath))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {

		filePath := filepath.Join(staticPath, r.URL.Path)
		stat, err := os.Stat(filePath)

		if err == nil && !stat.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

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

	r := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8081", "http://localhost:3000", "http://procure.ujaotech.com", "https://procure.ujaotech.com"},
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

	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	auth0Audience := os.Getenv("AUTH0_AUDIENCE")

	authValidator, err := appMiddleware.NewValidator(auth0Domain, auth0Audience)
	if err != nil {
		log.Fatalf("Failed to create authentication validator: %v", err)
	}

	r.Route("/api", func(apiRouter chi.Router) {
		handlers.RegisterUserRoutes(apiRouter)
		apiRouter.Group(func(authRouter chi.Router) {
			authRouter.Use(appMiddleware.TokenMiddleware(authValidator))
			authRouter.Post("/requisitions", handlers.CreateRequisitionHandler)
			authRouter.Get("/requisitions", handlers.ListRequisitionsHandler)
			authRouter.Get("/requisitions/{id}", handlers.GetRequisitionHandler)
			authRouter.Post("/requisitions/{id}/action", handlers.HandleRequisitionAction)

			tenderHandler := handlers.NewTenderHandler(db)
			authRouter.Post("/tenders", tenderHandler.CreateTender)
			authRouter.Get("/tenders/{id}", tenderHandler.GetTenderByID)
			bidHandler := handlers.NewBidHandler(db)
			authRouter.Post("/tenders/{tenderId}/bids", bidHandler.CreateBid)
			authRouter.Get("/tenders/{tenderId}/bids", bidHandler.ListTenderBids)
			authRouter.Get("/my-bids", bidHandler.ListMyBids)
			authRouter.Get("/dashboard/requisition-stats", handlers.GetRequisitionStatsHandler)
			authRouter.Get("/dashboard/recent-requisitions", handlers.GetRecentRequisitionsHandler)
			authRouter.Get("/dashboard/live-tenders", handlers.GetLiveTendersHandler)
			authRouter.Get("/dashboard/creation-rate", handlers.GetCreationRateHandler)
			authRouter.Get("/dashboard/my-stats", handlers.GetMyRequisitionStatsHandler)
			authRouter.Get("/dashboard/my-recent-requisitions", handlers.GetMyRecentRequisitionsHandler)
			authRouter.Get("/dashboard/supplier", handlers.GetSupplierDashboardDataHandler)
		})
	})

	serveFrontend(r, "frontend/dist")
	log.Println("Server starting on :8081, serving API and Frontend")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
