CREATE TABLE IF NOT EXISTS exercises (
    id           UUID        PRIMARY KEY,
    trainer_id   UUID        NOT NULL REFERENCES trainers(id),
    name         TEXT        NOT NULL,
    muscle_group TEXT,
    equipment    TEXT,
    video_url    TEXT,
    notes        TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_exercises_trainer_id ON exercises(trainer_id);
