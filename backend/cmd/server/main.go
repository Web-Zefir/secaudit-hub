package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/config"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/database"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/handlers"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/repository"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/service"
)

func main() {
	cfg := config.Load()
	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	users := repository.NewUserRepo(db)
	projects := repository.NewProjectRepo(db)
	scans := repository.NewScanRepo(db)
	findings := repository.NewFindingRepo(db)

	router := handlers.NewRouter(handlers.Deps{
		Config:    cfg,
		DB:        db,
		Auth:      service.NewAuthService(users),
		Projects:  service.NewProjectService(projects),
		Scans:     service.NewScanService(projects, scans, findings),
		Dashboard: service.NewDashboardService(projects, findings),
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("secaudit api listening on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
