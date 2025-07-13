package database

import (
	"log"
	"os"
	"path/filepath"
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

		// Ensure the directory for the database file exists.
		dir := filepath.Dir(dbPath)
		if mkdirErr := os.MkdirAll(dir, 0755); mkdirErr != nil {
			log.Printf("FATAL: Failed to create database directory %s: %v\n", dir, mkdirErr)
			err = mkdirErr
			return
		}

		log.Println("Attempting to connect to SQLite database using GORM...")
		var gormErr error
		db, gormErr = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if gormErr != nil {
			log.Printf("FATAL: Failed to connect to database %s: %v\n", dbPath, gormErr)
			err = gormErr
			return
		}

		sqlDB, dbPoolErr := db.DB()
		if dbPoolErr != nil {
			log.Printf("FATAL: Failed to get underlying sql.DB for pool settings: %v\n", dbPoolErr)
			err = dbPoolErr
			return
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		log.Printf("Successfully connected to SQLite database '%s' using GORM!", dbPath)
	})
	return err
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
