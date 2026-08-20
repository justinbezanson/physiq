package middleware

import (
	"net/http"
	"physiq/backend/internal/session"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const userIDKey = "auth.userID"

func Session(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie(session.CookieName)
		if err != nil || cookieValue == "" {
			c.Next()
			return
		}

		tokenHash, ok := session.HashFromCookieValue(cookieValue)
		if !ok {
			c.Next()
			return
		}

		var userID int64
		err = pool.QueryRow(c.Request.Context(),
			"SELECT user_id FROM sessions WHERE token_hash = $1 AND expires_at > now()",
			tokenHash,
		).Scan(&userID)
		if err == nil {
			c.Set(userIDKey, userID)
		}
		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := UserIDFrom(c); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func UserIDFrom(c *gin.Context) (int64, bool) {
	v, ok := c.Get(userIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}
