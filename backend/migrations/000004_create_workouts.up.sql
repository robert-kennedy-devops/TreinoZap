CREATE TABLE IF NOT EXISTS workouts (
    id         UUID        PRIMARY KEY,
    trainer_id UUID        NOT NULL REFERENCES trainers(id),
    client_id  UUID        NOT NULL REFERENCES clients(id),
    name       TEXT        NOT NULL,
    status     TEXT        NOT NULL DEFAULT 'draft',
    starts_at  DATE,
    ends_at    DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_one_active_workout_per_client
    ON workouts(client_id) WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_workouts_trainer_id ON workouts(trainer_id);
CREATE INDEX IF NOT EXISTS idx_workouts_client_id  ON workouts(client_id);
CREATE INDEX IF NOT EXISTS idx_workouts_status     ON workouts(status);

CREATE TABLE IF NOT EXISTS workout_sections (
    id          UUID        PRIMARY KEY,
    workout_id  UUID        NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    description TEXT,
    order_index INT         NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workout_sections_workout_id ON workout_sections(workout_id);

CREATE TABLE IF NOT EXISTS workout_exercises (
    id             UUID        PRIMARY KEY,
    section_id     UUID        NOT NULL REFERENCES workout_sections(id) ON DELETE CASCADE,
    exercise_id    UUID        REFERENCES exercises(id),
    exercise_name  TEXT        NOT NULL,
    sets           TEXT,
    reps           TEXT,
    rest_seconds   INT,
    load_note      TEXT,
    technique_note TEXT,
    video_url      TEXT,
    order_index    INT         NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workout_exercises_section_id ON workout_exercises(section_id);
