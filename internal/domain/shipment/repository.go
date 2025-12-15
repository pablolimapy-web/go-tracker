// internal/domain/shipment/repository.go
package shipment

import "context"

type Repository interface {
	Create(ctx context.Context, s Shipment) (Shipment, error)
	FindByID(ctx context.Context, id int64) (Shipment, error)
	ListByUser(ctx context.Context, userID int64) ([]Shipment, error)
	UpdateStatus(ctx context.Context, id int64, status Status) error
}
