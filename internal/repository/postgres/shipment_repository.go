// internal/repository/postgres/shipment_repository.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pablolimapy-web/go-tracker/internal/domain/shipment"
)

type ShipmentRepository struct {
	db *sql.DB
}

func NewShipmentRepository(db *sql.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) Create(ctx context.Context, s shipment.Shipment) (shipment.Shipment, error) {
	const q = `
        INSERT INTO shipments (user_id, code, carrier, status, last_update_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;
    `

	err := r.db.QueryRowContext(ctx, q,
		s.UserID, s.Code, s.Carrier, s.Status, s.LastUpdateAt, s.CreatedAt,
	).Scan(&s.ID)
	if err != nil {
		return shipment.Shipment{}, fmt.Errorf("Create: %w", err)
	}

	return s, nil
}
