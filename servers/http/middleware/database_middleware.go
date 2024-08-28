package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mi0772/nuke-go/engine"
)

func DatabaseMiddleware(db *engine.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}

}
