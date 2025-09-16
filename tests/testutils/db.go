package testutils

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"rest_api/config"
	"rest_api/models"

	"gorm.io/gorm"
)

// ConfigureTestDB sets environment and connects GORM to the test Postgres DB.
func ConfigureTestDB() *gorm.DB {
	// Allow override via env in CI; otherwise default to local docker compose
	dsn := os.Getenv("DB_CONNECTION_STRING")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
		os.Setenv("DB_CONNECTION_STRING", dsn)
	}

	// Reuse existing connect logic from app
	config.ConnectToDB()
	// Ensure schema exists
	_ = config.DB.AutoMigrate(&models.Post{}, &models.User{})

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

// ResetDatabase truncates all data and loads the tests/seeds.sql file.
func ResetDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// Clean tables (restart identity to keep deterministic IDs)
	if _, err = sqlDB.Exec("TRUNCATE TABLE posts, users RESTART IDENTITY CASCADE;"); err != nil {
		// Tables might not exist yet, try migrate then truncate again
		_ = db.AutoMigrate(&models.Post{}, &models.User{})
		_, _ = sqlDB.Exec("TRUNCATE TABLE posts, users RESTART IDENTITY CASCADE;")
	}

	// Load embedded seeds.sql
	statements := splitSQLStatements(seedSQL)
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := sqlDB.Exec(stmt); err != nil {
			return fmt.Errorf("seed exec error: %w (stmt: %s)", err, stmt)
		}
	}
	return nil
}

// splitSQLStatements is a naive splitter by semicolon that handles simple files.
func splitSQLStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

//go:embed seeds.sql
var seedSQL string
