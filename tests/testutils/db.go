package testutils

import (
	"context"
	"os"
	"time"

	"rest_api/config"
	"rest_api/models"

	"gorm.io/gorm"
)

// ConfigureTestDB sets environment and connects GORM to the test Postgres DB.
func ConfigureTestDB() *gorm.DB {
	// Rely on environment variables; provide a sane default for test DB
	if os.Getenv("DB_CONNECTION_STRING") == "" {
		_ = os.Setenv("DB_CONNECTION_STRING", "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=UTC")
	}

	dsn := os.Getenv("DB_CONNECTION_STRING")
	// Reuse existing connect logic from app
	config.ConnectToDB()
	// Ensure schema exists
	_ = config.DB.AutoMigrate(&models.User{}, &models.Post{})

	waitForPostgres(dsn)

	return config.DB
}

// waitForPostgres tries simple connections until the DB responds or times out.
func waitForPostgres(dsn string) {
	// GORM uses pgx; for low-level readiness use database/sql with pgx too
	// But we can simply retry gorm ping via underlying sql.DB
	sqlDB, err := config.DB.DB()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		if e := sqlDB.PingContext(ctx); e == nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// BeginTxWithSeeds starts a DB transaction and swaps the global config.DB to
// point to this transaction for the duration of a test. No seeds are loaded.
// The returned cleanup must be deferred to rollback the transaction and
// restore the original global DB.
func BeginTxWithSeeds() (func(), error) {
	db := ConfigureTestDB()
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	original := config.DB
	config.DB = tx

	cleanup := func() {
		_ = tx.Rollback()
		config.DB = original
	}

	return cleanup, nil
}

// ResetDatabase is a no-op retained for backwards compatibility.
func ResetDatabase(db *gorm.DB) error { return nil }
