package handlers

import (
	"errors"
	"net/http"
	"physiq/backend/internal/config"
	"physiq/backend/internal/middleware"
	"physiq/backend/internal/session"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("physiq-timing-equalizer"), 10)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		req.Email = strings.TrimSpace(req.Email)

		var id int64
		var name, storedHash string

		err := pool.QueryRow(
			c.Request.Context(),
			"SELECT id, name, password_hash FROM users WHERE email = $1", req.Email,
		).Scan(&id, &name, &storedHash)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not login"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		cookieValue, tokenHash, err := session.Generate()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not login"})
			return
		}

		expiresAt := time.Now().Add(session.TTL)
		_, err = pool.Exec(c.Request.Context(),
			`INSERT INTO sessions (user_id, token_hash, ip_address, user_agent, expires_at)
			 VALUES ($1, $2, $3, $4, $5)`,
			id, tokenHash, c.ClientIP(), c.Request.UserAgent(), expiresAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not login"})
			return
		}

		_, _ = pool.Exec(c.Request.Context(),
			"DELETE FROM sessions WHERE expires_at < now()")

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     session.CookieName,
			Value:    cookieValue,
			Path:     "/",
			HttpOnly: true,
			Secure:   cfg.AppMode != "DEV",
			SameSite: http.SameSiteLaxMode,
			Expires:  expiresAt,
		})

		c.JSON(http.StatusOK, gin.H{"id": id, "name": name})
	}
}

func Logout(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookieValue, err := c.Cookie(session.CookieName); err == nil && cookieValue != "" {
			if tokenHash, ok := session.HashFromCookieValue(cookieValue); ok {
				_, _ = pool.Exec(c.Request.Context(),
					"DELETE FROM sessions WHERE token_hash = $1", tokenHash)
			}
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     session.CookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   cfg.AppMode != "DEV",
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func Me(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := middleware.UserIDFrom(c)

		var name, email string
		err := pool.QueryRow(c.Request.Context(),
			"SELECT name, email FROM users WHERE id = $1", userID,
		).Scan(&name, &email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": userID, "name": name, "email": email})
	}
}

func Register(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
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
			`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
			req.Name, req.Email, string(hash),
		).Scan(&id)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
				return
			}
			if cfg.AppMode == "DEV" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"id": id, "email": req.Email, "name": req.Name})
	}
}
