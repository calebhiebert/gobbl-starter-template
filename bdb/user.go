package bdb

import (
	"context"
	"time"

	"github.com/avast/retry-go"
	"upper.io/db.v3/lib/sqlbuilder"
)

// User represents the holiday_type_interests db table
type User struct {
	ID        string    `db:"id,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`

	CountryCode         *string `db:"country_code"`
	CountryCodeOriginal *string `db:"country_code_original"`
	FirstName           *string `db:"first_name"`
	LastName            *string `db:"last_name"`
	Email               *string `db:"email"`
}

// DeleteUserData will delete all of a user's data from the database
// TODO modify this method with any new tables that contain user data
func (d *DBSession) DeleteUserData(id string) error {
	ctx := context.Background()

	err := retry.Do(func() error {

		// Perform user deletion inside a transaction
		return d.Tx(ctx, func(tx sqlbuilder.Tx) error {
			_, err := tx.DeleteFrom("user_actions").
				Where("user_id = ?", id).
				Exec()
			if err != nil {
				return err
			}

			_, err = tx.DeleteFrom("users").
				Where("id = ?", id).
				Exec()
			if err != nil {
				return err
			}

			return nil
		})
	})
	if err != nil {
		return err
	}

	return nil
}
