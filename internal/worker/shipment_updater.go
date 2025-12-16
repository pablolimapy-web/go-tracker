package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pablolimapy-web/go-tracker/internal/domain/shipment"
)

type ShipmentStatusProvider interface {
	NextStatus(ctx context.Context, s shipment.Shipment) (shipment.Status, error)
}

// Provider mock: PENDING -> IN_TRANSIT -> DELIVERED
type MockProvider struct{}

func (p *MockProvider) NextStatus(ctx context.Context, s shipment.Shipment) (shipment.Status, error) {
	switch s.Status {
	case shipment.StatusPending:
		return shipment.StatusInTransit, nil
	case shipment.StatusInTransit:
		return shipment.StatusDelivered, nil
	default:
		return s.Status, nil
	}
}

type ShipmentUpdater struct {
	repo        shipment.Repository
	provider    ShipmentStatusProvider
	interval    time.Duration
	batchSize   int
	concurrency int
}

func NewShipmentUpdater(
	repo shipment.Repository,
	provider ShipmentStatusProvider,
	interval time.Duration,
	batchSize int,
	concurrency int,
) *ShipmentUpdater {
	if batchSize <= 0 {
		batchSize = 50
	}
	if concurrency <= 0 {
		concurrency = 4
	}
	if interval <= 0 {
		interval = 10 * time.Second
	}

	return &ShipmentUpdater{
		repo:        repo,
		provider:    provider,
		interval:    interval,
		batchSize:   batchSize,
		concurrency: concurrency,
	}
}

func (u *ShipmentUpdater) Run(ctx context.Context) {
	ticker := time.NewTicker(u.interval)
	defer ticker.Stop()

	log.Println("worker: shipment updater started")

	// roda já na largada
	u.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("worker: shipment updater stopped")
			return
		case <-ticker.C:
			u.tick(ctx)
		}
	}
}

func (u *ShipmentUpdater) tick(ctx context.Context) {
	shipments, err := u.repo.ListPending(ctx, u.batchSize)
	if err != nil {
		log.Printf("worker: ListPending error: %v", err)
		return
	}
	if len(shipments) == 0 {
		return
	}

	sem := make(chan struct{}, u.concurrency)
	var wg sync.WaitGroup

	for _, s := range shipments {
		s := s

		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			next, err := u.provider.NextStatus(ctx, s)
			if err != nil {
				log.Printf("worker: NextStatus id=%d error: %v", s.ID, err)
				return
			}

			if next == s.Status {
				return
			}

			if err := u.repo.UpdateStatus(ctx, s.ID, next); err != nil {
				log.Printf("worker: UpdateStatus id=%d error: %v", s.ID, err)
				return
			}

			log.Printf("worker: shipment id=%d status %s -> %s", s.ID, s.Status, next)
		}()
	}

	wg.Wait()
}
