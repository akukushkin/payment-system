CREATE TABLE IF NOT EXISTS wallet(
    id BIGSERIAL PRIMARY KEY,
    value BIGINT NOT NULL DEFAULT 0,
    idempotency_key VARCHAR(36) UNIQUE,
    CONSTRAINT value_non_negative CHECK (value >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS wallet_idempotency_key_unique_idx
    ON wallet(idempotency_key);

CREATE TABLE IF NOT EXISTS operation(
    wallet_id BIGINT NOT NULL,
    value BIGINT NOT NULL,
    direction SMALLINT NOT NULL,
    date DATE DEFAULT now()::DATE,
    idempotency_key VARCHAR(36) NOT NULL,
    CONSTRAINT fk_wallet FOREIGN KEY(wallet_id) REFERENCES wallet(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS operation_idempotency_key_wallet_id_unique_idx
    ON operation(idempotency_key, wallet_id);
