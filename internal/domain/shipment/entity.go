// internal/domain/shipment/entity.go
package shipment

import "time"

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusError     Status = "ERROR"
)

type Shipment struct {
	ID           int64
	UserID       int64
	Code         string
	Carrier      string
	Status       Status
	LastUpdateAt time.Time
	CreatedAt    time.Time
}
