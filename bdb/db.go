package bdb

import (
	"os"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

// DBSession is a wrapper around upper.db's database struct
// this allows for easier custom db methods to be added
type DBSession struct {
	sqlbuilder.Database
}

// New creates a new database session
func New() (*DBSession, error) {
	sess, err := postgresql.Open(getDBSettings())
	if err != nil {
		return nil, err
	}

	return &DBSession{
		Database: sess,
	}, nil
}

func getDBSettings() postgresql.ConnectionURL {
	settings := postgresql.ConnectionURL{
		Host:     "localhost:5432",
		Database: "bdb",
		User:     "bdb",
		Password: "bdb",
	}

	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")

	if host != "" {
		settings.Host = host
	}

	if dbName != "" {
		settings.Database = dbName
	}

	if dbPass != "" {
		settings.Password = dbPass
	}

	if username != "" {
		settings.User = username
	}

	if os.Getenv("DB_SSL") == "true" {
		settings.Options = map[string]string{
			"sslmode": "require",
		}
	}

	return settings
}
