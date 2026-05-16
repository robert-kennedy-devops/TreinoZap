CREATE TABLE IF NOT EXISTS whatsapp_channels (
    id                UUID        PRIMARY KEY,
    name              TEXT        NOT NULL,
    phone             TEXT,
    jid               TEXT,
    status            TEXT        NOT NULL DEFAULT 'disconnected',
    is_default        BOOLEAN     NOT NULL DEFAULT false,
    last_connected_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS whatsapp_messages (
    id                  UUID        PRIMARY KEY,
    channel_id          UUID        REFERENCES whatsapp_channels(id),
    trainer_id          UUID        REFERENCES trainers(id),
    client_id           UUID        REFERENCES clients(id),
    direction           TEXT        NOT NULL CHECK (direction IN ('inbound', 'outbound')),
    phone               TEXT        NOT NULL,
    message             TEXT        NOT NULL,
    command             TEXT,
    status              TEXT,
    provider_message_id TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_trainer_id ON whatsapp_messages(trainer_id);
CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_client_id  ON whatsapp_messages(client_id);
CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_phone      ON whatsapp_messages(phone);
CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_direction  ON whatsapp_messages(direction);
