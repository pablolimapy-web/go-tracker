// internal/domain/shipment/service.go
package shipment

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateShipment(
	ctx context.Context,
	userID int64,
	code, carrier string,
) (Shipment, error) {
	sh := Shipment{
		UserID:       userID,
		Code:         code,
		Carrier:      carrier,
		Status:       StatusPending,
		LastUpdateAt: time.Now(),
		CreatedAt:    time.Now(),
	}

	return s.repo.Create(ctx, sh)
}

func (s *Service) GetByID(ctx context.Context, id int64) (Shipment, error) {
	if id <= 0 {
		return Shipment{}, ErrInvalidInput
	}
	return s.repo.FindByID(ctx, id)
}
