package workout

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("treino não encontrado")
var ErrActiveWorkoutExists = errors.New("já existe um treino ativo para este cliente")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, w *Workout) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO workouts (id, trainer_id, client_id, name, status, starts_at, ends_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		w.ID, w.TrainerID, w.ClientID, w.Name, w.Status, w.StartsAt, w.EndsAt, w.CreatedAt, w.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, s := range w.Sections {
		_, err = tx.Exec(ctx, `
			INSERT INTO workout_sections (id, workout_id, name, description, order_index, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			s.ID, w.ID, s.Name, s.Description, s.OrderIndex, s.CreatedAt, s.UpdatedAt,
		)
		if err != nil {
			return err
		}

		for _, e := range s.Exercises {
			_, err = tx.Exec(ctx, `
				INSERT INTO workout_exercises
				(id, section_id, exercise_id, exercise_name, sets, reps, rest_seconds, load_note, technique_note, video_url, order_index, created_at, updated_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
				e.ID, s.ID, e.ExerciseID, e.ExerciseName, e.Sets, e.Reps, e.RestSeconds,
				e.LoadNote, e.TechniqueNote, e.VideoURL, e.OrderIndex, e.CreatedAt, e.UpdatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) FindByID(ctx context.Context, id, trainerID uuid.UUID) (*Workout, error) {
	query := `
		SELECT id, trainer_id, client_id, name, status,
		       to_char(starts_at, 'YYYY-MM-DD'), to_char(ends_at, 'YYYY-MM-DD'),
		       created_at, updated_at
		FROM workouts WHERE id=$1 AND trainer_id=$2`

	var w Workout
	err := r.db.QueryRow(ctx, query, id, trainerID).Scan(
		&w.ID, &w.TrainerID, &w.ClientID, &w.Name, &w.Status,
		&w.StartsAt, &w.EndsAt, &w.CreatedAt, &w.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	sections, err := r.findSections(ctx, w.ID)
	if err != nil {
		return nil, err
	}
	w.Sections = sections
	return &w, nil
}

func (r *Repository) FindActiveByClientID(ctx context.Context, clientID uuid.UUID) (*Workout, error) {
	query := `
		SELECT id, trainer_id, client_id, name, status,
		       to_char(starts_at, 'YYYY-MM-DD'), to_char(ends_at, 'YYYY-MM-DD'),
		       created_at, updated_at
		FROM workouts WHERE client_id=$1 AND status='active'`

	var w Workout
	err := r.db.QueryRow(ctx, query, clientID).Scan(
		&w.ID, &w.TrainerID, &w.ClientID, &w.Name, &w.Status,
		&w.StartsAt, &w.EndsAt, &w.CreatedAt, &w.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	sections, err := r.findSections(ctx, w.ID)
	if err != nil {
		return nil, err
	}
	w.Sections = sections
	return &w, nil
}

func (r *Repository) findSections(ctx context.Context, workoutID uuid.UUID) ([]Section, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, workout_id, name, description, order_index, created_at, updated_at
		FROM workout_sections WHERE workout_id=$1 ORDER BY order_index ASC`,
		workoutID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []Section
	for rows.Next() {
		var s Section
		if err := rows.Scan(&s.ID, &s.WorkoutID, &s.Name, &s.Description, &s.OrderIndex, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		exercises, err := r.findExercises(ctx, s.ID)
		if err != nil {
			return nil, err
		}
		s.Exercises = exercises
		sections = append(sections, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sections, nil
}

func (r *Repository) findExercises(ctx context.Context, sectionID uuid.UUID) ([]WorkoutExercise, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, section_id, exercise_id, exercise_name, sets, reps, rest_seconds,
		       load_note, technique_note, video_url, order_index, created_at, updated_at
		FROM workout_exercises WHERE section_id=$1 ORDER BY order_index ASC`,
		sectionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []WorkoutExercise
	for rows.Next() {
		var e WorkoutExercise
		if err := rows.Scan(
			&e.ID, &e.SectionID, &e.ExerciseID, &e.ExerciseName, &e.Sets, &e.Reps, &e.RestSeconds,
			&e.LoadNote, &e.TechniqueNote, &e.VideoURL, &e.OrderIndex, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (r *Repository) ListByClient(ctx context.Context, clientID, trainerID uuid.UUID) ([]Workout, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, trainer_id, client_id, name, status,
		       to_char(starts_at, 'YYYY-MM-DD'), to_char(ends_at, 'YYYY-MM-DD'),
		       created_at, updated_at
		FROM workouts WHERE client_id=$1 AND trainer_id=$2 ORDER BY created_at DESC`,
		clientID, trainerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []Workout
	for rows.Next() {
		var w Workout
		if err := rows.Scan(
			&w.ID, &w.TrainerID, &w.ClientID, &w.Name, &w.Status,
			&w.StartsAt, &w.EndsAt, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *Repository) Activate(ctx context.Context, id, clientID, trainerID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Archive previously active workout
	_, err = tx.Exec(ctx, `
		UPDATE workouts SET status='archived', updated_at=$1
		WHERE client_id=$2 AND status='active' AND id!=$3`,
		time.Now().UTC(), clientID, id,
	)
	if err != nil {
		return err
	}

	res, err := tx.Exec(ctx, `
		UPDATE workouts SET status='active', updated_at=$1
		WHERE id=$2 AND trainer_id=$3`,
		time.Now().UTC(), id, trainerID,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return tx.Commit(ctx)
}

func (r *Repository) UpdateWithSections(ctx context.Context, w *Workout) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx, `
		UPDATE workouts SET name=$1, status=$2, starts_at=$3, ends_at=$4, updated_at=$5
		WHERE id=$6 AND trainer_id=$7`,
		w.Name, w.Status, w.StartsAt, w.EndsAt, w.UpdatedAt, w.ID, w.TrainerID,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	// Delete old sections (cascades to exercises)
	_, err = tx.Exec(ctx, `DELETE FROM workout_sections WHERE workout_id=$1`, w.ID)
	if err != nil {
		return err
	}

	for _, s := range w.Sections {
		_, err = tx.Exec(ctx, `
			INSERT INTO workout_sections (id, workout_id, name, description, order_index, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			s.ID, w.ID, s.Name, s.Description, s.OrderIndex, s.CreatedAt, s.UpdatedAt,
		)
		if err != nil {
			return err
		}
		for _, e := range s.Exercises {
			_, err = tx.Exec(ctx, `
				INSERT INTO workout_exercises
				(id, section_id, exercise_id, exercise_name, sets, reps, rest_seconds, load_note, technique_note, video_url, order_index, created_at, updated_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
				e.ID, s.ID, e.ExerciseID, e.ExerciseName, e.Sets, e.Reps, e.RestSeconds,
				e.LoadNote, e.TechniqueNote, e.VideoURL, e.OrderIndex, e.CreatedAt, e.UpdatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) Archive(ctx context.Context, id, trainerID uuid.UUID) error {
	res, err := r.db.Exec(ctx, `
		UPDATE workouts SET status='archived', updated_at=$1 WHERE id=$2 AND trainer_id=$3`,
		time.Now().UTC(), id, trainerID,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
