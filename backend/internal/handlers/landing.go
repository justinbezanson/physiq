package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"physiq/backend/templates"
)

func Landing() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Status(http.StatusOK)
		if err := templates.Index().Render(c.Request.Context(), c.Writer); err != nil {
			_ = c.Error(err)
		}
	}
}
