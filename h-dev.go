package main

import (
	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl/messenger"
	"gobbl-template/bdb"
)

// TCID is what the user should enter to trigger the HGetID handler
var TCID = "GET_ID"

// TCDeleteData is what the user should enter to trigger the HDeleteUserData handler
var TCDeleteData = "DELETE_DATA"

// HGetID will respond with the user's ID
func HGetID(c *gbl.Context) {
	r := fb.CreateResponse(c)
	r.Text(c.User.ID)
}

// HDeleteUserData will delete all of a user's data
func HDeleteUserData(c *gbl.Context) {
	err := bdb.DB(c).DeleteUserData(c.User.ID)
	if err != nil {
		c.Errorf("Error occured %v", err)
	}
}
