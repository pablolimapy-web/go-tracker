package main

import (
	"log"
	"net/http"

	dbpostgres "github.com/pablolimapy-web/go-tracker/internal/database/postgres"
	repopostgres "github.com/pablolimapy-web/go-tracker/internal/repository/postgres"

	"github.com/pablolimapy-web/go-tracker/internal/domain/shipment"
	"github.com/pablolimapy-web/go-tracker/internal/http/handler"
	"github.com/pablolimapy-web/go-tracker/internal/http/router"
)

func main() {
	dsn := "postgres://tracker:tracker@localhost:5434/tracker?sslmode=disable"

	db := dbpostgres.NewConnection(dsn)

	repo := repopostgres.NewShipmentRepository(db)
	service := shipment.NewService(repo)
	h := handler.NewShipmentHandler(service)

	r := router.New(h)

	log.Println("API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
