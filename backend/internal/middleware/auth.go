package middleware

import (
	"github.com/gin-gonic/gin"
)

// Session resolves the current user from the session cookie.
// TODO: wire up real session lookup once auth lands.
func Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
