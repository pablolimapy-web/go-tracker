// internal/worker/updater.go
package worker

import (
	"context"
	"log"
	"time"

	"github.com/pablolimapy-web/go-tracker/internal/domain/shipment"
)

type StatusProvider interface {
	FetchStatus(ctx context.Context, code, carrier string) (shipment.Status, error)
}

type Updater struct {
	repo     shipment.Repository
	provider StatusProvider
}

func NewUpdater(r shipment.Repository, p StatusProvider) *Updater {
	return &Updater{repo: r, provider: p}
}

func (u *Updater) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("updater stopped")
			return
		case <-ticker.C:
			u.tick(ctx)
		}
	}
}

func (u *Updater) tick(ctx context.Context) {
	// aqui você buscaria shipments pendentes e atualizaria um por um
	// simulando:
	log.Println("Running updater tick...")
}
