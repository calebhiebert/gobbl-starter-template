package bdb

import (
	"github.com/calebhiebert/gobbl"
	"github.com/gin-gonic/gin"
)

// Middleware is a middleware function that adds the database session to the gobbl context.
// This is done to make using the database in handler functions easier
func Middleware(dbSession *DBSession) gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		// Set the database session on the context so it can be retrieved by the helper function later
		c.Flag("__dbsession", dbSession)
		c.Next()
	}
}

// GinMiddleware creates a gin middleware that will store the database session
// on incoming requests
func GinMiddleware(dbSession *DBSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("__dbsession", dbSession)
		c.Next()
	}
}

// DB will grab the db session from the gobbl context and return it
func DB(c *gbl.Context) *DBSession {
	if c.HasFlag("__dbsession") {
		return c.GetFlag("__dbsession").(*DBSession)
	}

	panic("Database was not found on context!")
}

// GDB will retrieve the database session from the gin context
// and cast/return it for us
func GDB(c *gin.Context) *DBSession {
	db, exists := c.Get("__dbsession")
	if exists == true {
		return db.(*DBSession)
	}

	panic("Database was not found on gin context!")
}
