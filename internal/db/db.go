package db

import (
	"github.com/kgantsov/filo/internal/model"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the SQLite database and runs auto-migration for the File struct.
func InitDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal().Msgf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&model.File{})
	if err != nil {
		log.Fatal().Msgf("Failed to auto-migrate database: %v", err)
	}

	log.Info().Msg("Database connection established and migrated.")
	return db
}
