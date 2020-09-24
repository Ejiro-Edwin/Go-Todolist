package database

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

//Connect creates a new database connection
func Connect() (*sqlx.DB, error) {
	//Connect to database:
	logrus.Debug("Connecting to database.")
	conn, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to database")
	}

	conn.SetMaxOpenConns(32)

	//Check if database running
	if err := waitForDb(conn.DB); err != nil {
		return nil, err
	}

	//Migrate database schema
	if err := migrateDb(conn.DB); err != nil {
		return nil, errors.Wrap(err, "could not migrate database")
	}

	return conn, nil
}

//New creates a new database
func New() (Database, error) {
	conn, err := Connect()
	if err != nil {
		return nil, err
	}

	d := &database{
		conn: conn,
	}

	return d, nil
}

func waitForDb(conn *sql.DB) error {
	ready := make(chan struct{})
	go func() {
		for {
			if err := conn.Ping(); err == nil {
				close(ready)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return nil
	case <-time.After(time.Duration(5000) * time.Millisecond):
		return errors.New("database not ready")
	}
}

