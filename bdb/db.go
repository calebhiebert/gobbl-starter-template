package bdb

import (
	"fmt"
	"os"
	"time"

	retry "github.com/avast/retry-go"
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
	var dbSess sqlbuilder.Database

	err := retry.Do(
		func() error {
			sess, err := postgresql.Open(getDBSettings())
			if err != nil {
				return err
			}

			dbSess = sess
			return nil
		},
		retry.Delay(1*time.Second),
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("retrying database connection %d: %v\n", n, err)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &DBSession{
		Database: dbSess,
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
