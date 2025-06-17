package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the postgres driver

	"github.com/kylenguyen/wallet-app/internal/config"
)

// Connect establishes a connection to a PostgreSQL database.
// It configures the connection using parameters from the application's configuration.
//
// Parameters:
//   - config: The database configuration.
//
// Returns:
//   - *sqlx.DB: A pointer to the connected database instance.
//   - error: An error if the connection fails.
func Connect(config config.Config) (*sqlx.DB, error) {
	// DSN for PostgreSQL
	// Example: "host=localhost port=5432 user=user password=password dbname=mydb sslmode=disable"
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseVar.Host,
		config.DatabaseVar.Port,
		config.DatabaseVar.User,
		config.DatabaseVar.Password,
		config.DatabaseVar.Name,
	)

	// Use sqlx.Open for a standard connection without DataDog tracing
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		// It's generally better to return the error rather than panic,
		// allowing the caller to decide how to handle it.
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(config.DatabaseVar.MaxOpenConns)
	db.SetMaxIdleConns(config.DatabaseVar.MaxIdleConns)
	db.SetConnMaxLifetime(config.DatabaseVar.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		// Close the connection if ping fails to prevent resource leaks
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
