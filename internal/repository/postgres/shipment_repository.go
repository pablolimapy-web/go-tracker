package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pablolimapy-web/go-tracker/internal/domain/shipment"
)

type ShipmentRepository struct {
	db *sql.DB
}

func NewShipmentRepository(db *sql.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) Create(
	ctx context.Context,
	s shipment.Shipment,
) (shipment.Shipment, error) {

	const q = `
        INSERT INTO shipments (
            user_id, code, carrier, status, last_update_at, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;
    `

	err := r.db.QueryRowContext(ctx, q,
		s.UserID,
		s.Code,
		s.Carrier,
		s.Status,
		s.LastUpdateAt,
		s.CreatedAt,
	).Scan(&s.ID)

	if err != nil {
		return shipment.Shipment{}, fmt.Errorf("shipment create: %w", err)
	}

	return s, nil
}
func (r *ShipmentRepository) FindByID(ctx context.Context, id int64) (shipment.Shipment, error) {
	const q = `
        SELECT id, user_id, code, carrier, status, last_update_at, created_at
        FROM shipments
        WHERE id = $1
    `

	var s shipment.Shipment
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&s.ID,
		&s.UserID,
		&s.Code,
		&s.Carrier,
		&s.Status,
		&s.LastUpdateAt,
		&s.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shipment.Shipment{}, shipment.ErrNotFound
		}
		return shipment.Shipment{}, fmt.Errorf("FindByID: %w", err)
	}

	return s, nil
}
func (r *ShipmentRepository) ListByUser(ctx context.Context, userID int64) ([]shipment.Shipment, error) {
	const q = `
        SELECT id, user_id, code, carrier, status, last_update_at, created_at
        FROM shipments
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("ListByUser query: %w", err)
	}
	defer rows.Close()

	shipments := make([]shipment.Shipment, 0)

	for rows.Next() {
		var s shipment.Shipment
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Code,
			&s.Carrier,
			&s.Status,
			&s.LastUpdateAt,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListByUser scan: %w", err)
		}

		shipments = append(shipments, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListByUser rows: %w", err)
	}

	return shipments, nil
}

func (r *ShipmentRepository) UpdateStatus(ctx context.Context, id int64, status shipment.Status) error {
	const q = `
        UPDATE shipments
        SET status = $1,
            last_update_at = NOW(),
            locked_until = NULL,
            locked_by = NULL
        WHERE id = $2
    `

	res, err := r.db.ExecContext(ctx, q, status, id)
	if err != nil {
		return fmt.Errorf("UpdateStatus exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateStatus rowsAffected: %w", err)
	}

	if rows == 0 {
		return shipment.ErrNotFound
	}

	return nil
}
func (r *ShipmentRepository) ClaimPending(
	ctx context.Context,
	limit int,
	lockFor time.Duration,
	workerID string,
) ([]shipment.Shipment, error) {
	if limit <= 0 {
		limit = 50
	}
	if lockFor <= 0 {
		lockFor = 45 * time.Second
	}
	if workerID == "" {
		workerID = "worker"
	}

	// Claim atômico:
	// 1) escolhe jobs pendentes cujo lock expirou
	// 2) FOR UPDATE SKIP LOCKED impede outra instância de pegar os mesmos
	// 3) UPDATE marca locked_until/locked_by e retorna as linhas
	const q = `
        WITH picked AS (
            SELECT id
            FROM shipments
            WHERE status IN ('PENDING', 'IN_TRANSIT')
              AND (locked_until IS NULL OR locked_until < NOW())
            ORDER BY last_update_at ASC
            LIMIT $1
            FOR UPDATE SKIP LOCKED
        )
        UPDATE shipments s
        SET locked_until = NOW() + ($2 || ' seconds')::interval,
            locked_by = $3
        FROM picked
        WHERE s.id = picked.id
        RETURNING s.id, s.user_id, s.code, s.carrier, s.status, s.last_update_at, s.created_at;
    `

	rows, err := r.db.QueryContext(ctx, q, limit, int(lockFor.Seconds()), workerID)
	if err != nil {
		return nil, fmt.Errorf("ClaimPending query: %w", err)
	}
	defer rows.Close()

	out := make([]shipment.Shipment, 0, limit)
	for rows.Next() {
		var s shipment.Shipment
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Code,
			&s.Carrier,
			&s.Status,
			&s.LastUpdateAt,
			&s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("ClaimPending scan: %w", err)
		}
		out = append(out, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ClaimPending rows: %w", err)
	}

	return out, nil
}
