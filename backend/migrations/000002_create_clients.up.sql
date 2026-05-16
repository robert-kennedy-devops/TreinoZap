CREATE TABLE IF NOT EXISTS clients (
    id         UUID        PRIMARY KEY,
    trainer_id UUID        NOT NULL REFERENCES trainers(id),
    name       TEXT        NOT NULL,
    phone      TEXT        NOT NULL UNIQUE,
    status     TEXT        NOT NULL DEFAULT 'active',
    goal       TEXT,
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_clients_trainer_id ON clients(trainer_id);
CREATE INDEX IF NOT EXISTS idx_clients_phone ON clients(phone);
