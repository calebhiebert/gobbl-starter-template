package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gobbl-template/bdb"
	"gobbl-template/linkr"
)

// GHLinkr is a gin handler that can process linkr links
func GHLinkr(c *gin.Context) {
	if c.Query("d") == "" {
		c.JSON(400, gin.H{"error": "missing query param \"d\""})
		return
	}

	decodedLink, err := linkr.Decode(c.Query("d"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid linkr data param", "ctx": err})
		return
	}

	// Save to user action
	db := bdb.GDB(c)

	err = db.CreateUserActionURL(decodedLink.PSID, decodedLink.RedirectTo)
	if err != nil {
		fmt.Println("LINKR DB ERR TO USER ACTIONS", err)
	}

	c.Redirect(303, decodedLink.RedirectTo)
}
