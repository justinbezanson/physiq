package httpapi

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"physiq/backend/internal/config"
	"physiq/backend/internal/handlers"
	"physiq/backend/internal/middleware"
)

func NewRouter(cfg *config.Config, pool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.Health(pool))
	r.GET("/", handlers.Landing())

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register(pool, cfg))
	}

	api := r.Group("/api", middleware.Session())
	{
		api.GET("/health", handlers.APIHealth(pool))
	}

	r.Static("/static", cfg.StaticDir)
	r.GET("/app", spaHandler(cfg.StaticDir))

	r.NoRoute(spaFallback(cfg.StaticDir))

	return r
}

func spaHandler(staticDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		serveSPA(c, staticDir)
	}
}

func spaFallback(staticDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !wantsHTML(c.Request) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		serveSPA(c, staticDir)
	}
}

func serveSPA(c *gin.Context, staticDir string) {
	index := filepath.Join(staticDir, "index.html")
	data, err := os.ReadFile(index)
	if err != nil {
		c.String(http.StatusNotFound, "SPA not built")
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(data)
}

func wantsHTML(r *http.Request) bool {
	for _, accept := range r.Header.Values("Accept") {
		if contains(accept, "text/html") {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
