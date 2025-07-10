package database

import (
	"log"
	"os"
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

func getDatabasePath() string {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./procurement.db"
	}
	return dbPath
}

func InitDB() error {
	var err error
	once.Do(func() {
		dbPath := getDatabasePath()
		log.Println("Attempting to connect to SQLite database using GORM...")
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Printf("FATAL: Failed to connect to database %s: %v\n", dbPath, err)
			return
		}

		sqlDB, DBPoolErr := db.DB()
		if DBPoolErr != nil {
			log.Printf("FATAL: Failed to get underlying sql.DB for pool settings: %v\n", DBPoolErr)
			err = DBPoolErr // Propagate this error
			return
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		log.Printf("Successfully connected to SQLite database '%s' using GORM!", dbPath)
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
