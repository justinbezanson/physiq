package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	StaticDir string
}

func Load() Config {
	return Config{
		Port:      getenv("PORT", "8080"),
		DBHost:    getenv("DB_HOST", "localhost"),
		DBPort:    getenv("DB_PORT", "5432"),
		DBUser:    getenv("DB_USER", "physiq"),
		DBPass:    getenv("DB_PASSWORD", "physiq"),
		DBName:    getenv("DB_NAME", "physiq"),
		StaticDir: getenv("STATIC_DIR", "static"),
	}
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
