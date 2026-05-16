CREATE TABLE IF NOT EXISTS automation_rules (
    id         UUID        PRIMARY KEY,
    trainer_id UUID        NOT NULL REFERENCES trainers(id),
    keyword    TEXT        NOT NULL,
    action     TEXT        NOT NULL,
    is_active  BOOLEAN     NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(trainer_id, keyword)
);

CREATE INDEX IF NOT EXISTS idx_automation_rules_trainer_id ON automation_rules(trainer_id);
CREATE INDEX IF NOT EXISTS idx_automation_rules_keyword    ON automation_rules(keyword);
