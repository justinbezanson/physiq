package main

import (
	"context"
	"log"
	"net/http"

	"physiq/backend/internal/config"
	"physiq/backend/internal/db"
	"physiq/backend/internal/httpapi"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(context.Background(), cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: httpapi.NewRouter(cfg, pool),
	}

	log.Printf("physiq backend listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
