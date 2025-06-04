package database

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"procurement/models" // Uncommented to enable AutoMigrate for models.User
)

var (
	db   *gorm.DB
	once sync.Once
)

const DATABASE_NAME = "./procurement.db" // Or from env var

// InitDB initializes the GORM database connection.
// It's called once to set up the db pool.
func InitDB() error {
	var err error
	once.Do(func() {
		log.Println("Attempting to connect to SQLite database using GORM...")
		db, err = gorm.Open(sqlite.Open(DATABASE_NAME), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // Or logger.Silent for less noise
		})
		if err != nil {
			log.Printf("FATAL: Failed to connect to database %s: %v\n", DATABASE_NAME, err)
			// If we can't connect, we might want to os.Exit(1) or handle it more gracefully
			// For now, we return the error to be handled by the caller in main.go
			return
		}

		// Optional: Configure connection pool settings
		sqlDB, DBPoolErr := db.DB()
		if DBPoolErr != nil {
			log.Printf("FATAL: Failed to get underlying sql.DB for pool settings: %v\n", DBPoolErr)
			err = DBPoolErr // Propagate this error
			return
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		// sqlDB.SetConnMaxLifetime(time.Hour) // Example

		log.Printf("Successfully connected to SQLite database '%s' using GORM!", DATABASE_NAME)
	})
	return err // This err is from the gorm.Open or db.DB() call inside once.Do
}

// GetDB returns the initialized GORM database instance.
// Panics if InitDB has not been called or failed.
func GetDB() *gorm.DB {
	if db == nil {
		// This case should ideally not be reached if InitDB is called at startup
		// and its error is handled properly.
		log.Panicln("FATAL: GetDB called before InitDB or InitDB failed.")
	}
	return db
}

// SetupDatabaseSchema is now responsible for GORM auto-migration.
// It should be called after InitDB.
func SetupDatabaseSchema() {
	dbInstance := GetDB()
	if dbInstance == nil {
		log.Fatalf("FATAL: Cannot setup schema, database not initialized.")
		return
	}

	log.Println("INFO: Auto-migrating database schema with GORM...")

	err := dbInstance.AutoMigrate(
		&models.User{}, // Re-enabled AutoMigrate for User model
		&models.Requisition{},
		&models.RequisitionItem{},
		&models.Tender{}, // Add Tender model for auto-migration
		&models.Bid{},
		&models.BidItem{},
	)
	if err != nil {
		// If models.User was the only thing being migrated and it's commented out,
		// this error block might not be hit unless other models are added back.
		// However, if AutoMigrate([]) is called with no arguments and still returns an error,
		// this log will catch it.
		log.Fatalf("FATAL: Failed to auto-migrate database schema (or no models to migrate): %v", err)
	}

	log.Println("INFO: Database schema auto-migration complete (or no changes needed/skipped).")
}
