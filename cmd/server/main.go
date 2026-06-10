package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ervinas/server_uptime_checker/internal/api"
	"github.com/ervinas/server_uptime_checker/internal/checker"
	"github.com/ervinas/server_uptime_checker/internal/config"
	"github.com/ervinas/server_uptime_checker/internal/database"
	"github.com/ervinas/server_uptime_checker/internal/store"
	"github.com/ervinas/server_uptime_checker/internal/worker"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("database connect: %v", err)
	}
	defer db.Close()

	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "/app/migrations"
	}
	if err := database.Migrate(db, migrationsDir); err != nil {
		log.Fatalf("database migrate: %v", err)
	}

	s := store.New(db)
	chk := checker.New()
	w := worker.New(s, chk, cfg.CheckIntervalMinutes, cfg.CheckRetentionDays)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go w.Run(ctx)

	handler := api.NewHandler(s, chk)
	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           api.NewRouter(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)
}
