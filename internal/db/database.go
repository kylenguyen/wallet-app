package db

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	ddsql "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	ddsqlx "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/config"
)

// Connect establishes a connection to a MySQL database.
// It configures the connection using parameters from the application's configuration (Viper).
// It also integrates with DataDog for database tracing.
//
// Parameters:
//   - config: The database configuration.
//
// Returns:
//   - *sqlx.DB: A pointer to the connected database instance.
//   - error: An error if the connection fails.
func Connect(config config.Config) (*sqlx.DB, error) {
	dbServiceName := fmt.Sprintf("%s-%s-db", config.ServiceName, config.Env)

	ddsql.Register("mysql", &mysql.MySQLDriver{},
		ddsql.WithServiceName(dbServiceName),
		ddsql.WithAnalytics(true))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8&loc=Local",
		config.DatabaseVar.User,
		config.DatabaseVar.Password,
		config.DatabaseVar.Host,
		config.DatabaseVar.Port,
		config.DatabaseVar.Name)

	db, err := ddsqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.DatabaseVar.MaxOpenConns)

	db.SetMaxIdleConns(config.DatabaseVar.MaxIdleConns)

	db.SetConnMaxLifetime(config.DatabaseVar.ConnMaxLifetime)

	err = db.Ping()

	return db, err
}
