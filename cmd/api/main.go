package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	dbpostgres "github.com/pablolimapy-web/go-tracker/internal/database/postgres"
	repopostgres "github.com/pablolimapy-web/go-tracker/internal/repository/postgres"

	"github.com/pablolimapy-web/go-tracker/internal/domain/shipment"
	"github.com/pablolimapy-web/go-tracker/internal/http/handler"
	"github.com/pablolimapy-web/go-tracker/internal/http/middleware"
	"github.com/pablolimapy-web/go-tracker/internal/http/router"
	"github.com/pablolimapy-web/go-tracker/internal/worker"
)

func main() {
	dsn := "postgres://tracker:tracker@localhost:5434/tracker?sslmode=disable"

	db := dbpostgres.NewConnection(dsn)
	defer closeDB(db)

	repo := repopostgres.NewShipmentRepository(db)
	service := shipment.NewService(repo)
	h := handler.NewShipmentHandler(service)
	r := router.New(h)

	// Contexto "global" da aplicação (API + workers)
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// Worker: atualiza status periodicamente
	updater := worker.NewShipmentUpdater(
		repo,
		&worker.MockProvider{},
		10*time.Second, // intervalo
		50,             // batch
		4,              // concorrência
	)
	go updater.Run(appCtx)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.CORS(r),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("API listening on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Println("shutdown signal received")

	// 1) Para os workers
	appCancel()

	// 2) Para o servidor HTTP com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("http shutdown error: %v", err)
	} else {
		log.Println("http server stopped gracefully")
	}
}

func closeDB(db *sql.DB) {
	if db == nil {
		return
	}
	if err := db.Close(); err != nil {
		log.Printf("db close error: %v", err)
	} else {
		log.Println("db connection closed")
	}
}
