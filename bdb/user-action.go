package bdb

import (
	"time"

	"github.com/avast/retry-go"
)

// UserAction represents the user_actions table in the database
type UserAction struct {
	ID        uint      `db:"id,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`

	UserID     string  `db:"user_id"`
	URL        *string `db:"url,omitempty"`
	Button     *string `db:"button,omitempty"`
	QuickReply *string `db:"quick_reply,omitempty"`
	Message    *string `db:"message,omitempty"`
	Payload    *string `db:"payload,omitempty"`
}

// CreateUserActionButton will create a user action with the ButtonName
// property filled in
func (d *DBSession) CreateUserActionButton(psid, button, payload string) error {
	return retry.Do(func() error {
		_, err := d.InsertInto("user_actions").
			Columns("user_id", "button", "payload").
			Values(psid, button, payload).
			Exec()

		return err
	})
}

// CreateUserActionQR will create a user action with the QuickReply
// property filled in
func (d *DBSession) CreateUserActionQR(psid, qr, payload string) error {
	return retry.Do(func() error {
		_, err := d.InsertInto("user_actions").
			Columns("user_id", "quick_reply", "payload").
			Values(psid, qr, payload).
			Exec()

		return err
	})
}

// CreateUserActionMessage will create a user action with the Message
// property filled in
func (d *DBSession) CreateUserActionMessage(psid, message string) error {
	return retry.Do(func() error {
		_, err := d.InsertInto("user_actions").
			Columns("user_id", "message").
			Values(psid, message).
			Exec()

		return err
	})
}

// CreateUserActionURL will create a user action with the RedirectURL
// property filled in
func (d *DBSession) CreateUserActionURL(psid, url string) error {
	return retry.Do(func() error {
		_, err := d.InsertInto("user_actions").
			Columns("user_id", "url").
			Values(psid, url).
			Exec()

		return err
	})
}
