CREATE TABLE shipments (
                           id BIGSERIAL PRIMARY KEY,
                           user_id BIGINT NOT NULL,
                           code TEXT NOT NULL,
                           carrier TEXT NOT NULL,
                           status TEXT NOT NULL,
                           last_update_at TIMESTAMP NOT NULL,
                           created_at TIMESTAMP NOT NULL
);
