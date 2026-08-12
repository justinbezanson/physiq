package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	for _, key := range []string{"PORT", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "STATIC_DIR"} {
		t.Setenv(key, "")
	}

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "localhost")
	}
	if got := cfg.DatabaseURL(); got != "postgres://physiq:physiq@localhost:5432/physiq?sslmode=disable" {
		t.Errorf("DatabaseURL = %q", got)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DB_HOST", "db")

	cfg := Load()
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.DBHost != "db" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "db")
	}
}
