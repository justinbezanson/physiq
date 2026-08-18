package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		req.Email = strings.TrimSpace(req.Email)

		if req.Email == "" || !strings.Contains(req.Email, "@") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		}
		if len(req.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		var id int64
		err = pool.QueryRow(c.Request.Context(),
			`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
			req.Email, string(hash),
		).Scan(&id)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id, "email": req.Email})
	}
}
