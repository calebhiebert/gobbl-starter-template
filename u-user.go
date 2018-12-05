package main

import (
	"fmt"

	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl-starter-bot/bdb"
	"github.com/calebhiebert/gobbl/messenger"
	db "upper.io/db.v3"
)

// UUserExtractorMiddleware will return a middleware that
// grabs the psid out of the request and populates all data
// related to the user
func UUserExtractorMiddleware() gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		if c.User.ID == "" {
			c.Next()
			return
		}

		dbase := bdb.DB(c)

		var user bdb.User

		err := dbase.SelectFrom("users").
			Where("id = ?", c.User.ID).
			One(&user)
		if err != nil {
			if err == db.ErrNoMoreRows {
				// user does not exist
			} else {
				c.Log(10, fmt.Sprintf("Error querying database for user %v", err), "UserExtractor")
				c.Next()
				return
			}
		}

		// TODO check to see if the user is missing any information that we
		// can get from the facebook user info api, and update the db with
		// that information
		if user.ID == "" {
			HCreateFirstTimeUser(c)
			return
		}

		c.User.FirstName = *user.FirstName
		c.User.LastName = *user.LastName
		c.Flag("user", &user)
		c.Next()
	}
}

// HCreateFirstTimeUser will create a user in the database for the first time
func HCreateFirstTimeUser(c *gbl.Context) {
	mapi := c.Integration.(*fb.MessengerIntegration)

	userInfo, err := mapi.API.UserInfo(c.User.ID)
	if err != nil {
		c.Log(10, fmt.Sprintf("Error getting user data from fb %v", err), "UserExtractor")
		c.Next()
		return
	}

	infoMap := userInfo.(map[string]interface{})

	c.Infof("%+v", infoMap)

	firstName := infoMap["first_name"].(string)
	lastName := infoMap["last_name"].(string)

	c.User.FirstName = firstName
	c.User.LastName = lastName

	user := bdb.User{
		ID:        c.User.ID,
		FirstName: &firstName,
		LastName:  &lastName,
	}

	c.Flag("user", &user)

	_, err = bdb.DB(c).
		InsertInto("users").
		Columns("id", "first_name", "last_name").
		Values(c.User.ID, c.User.FirstName, c.User.LastName).
		Exec()
	if err != nil {
		c.Errorf("Erorr creating user %v", err)
	}

	c.Next()
}

// UUserActionLoggerMiddleware returns a middleware that will log
// all user actions in the database
func UUserActionLoggerMiddleware() gbl.MiddlewareFunction {
	return func(c *gbl.Context) {

		if c.HasFlag("fb:eventtype") {

			// process in a goroutine to avoid slowing down the rest of the request
			go func(c *gbl.Context) {
				var err error

				// Saving user actions in the database
				switch c.GetStringFlag("fb:eventtype") {
				case "quickreply":
					err = bdb.DB(c).CreateUserActionQR(c.User.ID, c.Request.Text, c.Request.Text)
				case "message":
					err = bdb.DB(c).CreateUserActionMessage(c.User.ID, c.Request.Text)
				case "payload":
					err = bdb.DB(c).CreateUserActionButton(c.User.ID, c.Request.Text, c.Request.Text)
				}

				if err != nil {
					c.Errorf("Error saving user action %v", err)
				}
			}(c)
		}

		c.Next()
	}
}
