ALTER TABLE shipments
    ADD COLUMN locked_until TIMESTAMP NULL,
ADD COLUMN locked_by TEXT NULL;

CREATE INDEX IF NOT EXISTS idx_shipments_lock
    ON shipments (status, locked_until);
