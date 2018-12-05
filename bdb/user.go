package bdb

import "time"

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
