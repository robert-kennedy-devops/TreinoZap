package exercise

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("exercício não encontrado")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, e *Exercise) error {
	query := `
		INSERT INTO exercises (id, trainer_id, name, muscle_group, equipment, video_url, notes, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.Exec(ctx, query,
		e.ID, e.TrainerID, e.Name, e.MuscleGroup, e.Equipment, e.VideoURL, e.Notes, e.CreatedAt, e.UpdatedAt,
	)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id, trainerID uuid.UUID) (*Exercise, error) {
	query := `
		SELECT id, trainer_id, name, muscle_group, equipment, video_url, notes, created_at, updated_at
		FROM exercises WHERE id=$1 AND trainer_id=$2`
	var e Exercise
	err := r.db.QueryRow(ctx, query, id, trainerID).Scan(
		&e.ID, &e.TrainerID, &e.Name, &e.MuscleGroup, &e.Equipment, &e.VideoURL, &e.Notes, &e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &e, err
}

type ListFilter struct {
	TrainerID uuid.UUID
	Search    string
	Page      int
	PageSize  int
}

func (r *Repository) List(ctx context.Context, f ListFilter) ([]Exercise, int, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}
	offset := (f.Page - 1) * f.PageSize

	where := "trainer_id=$1"
	args := []any{f.TrainerID}
	argN := 2

	if f.Search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR muscle_group ILIKE $%d OR equipment ILIKE $%d)", argN, argN+1, argN+2)
		like := "%" + f.Search + "%"
		args = append(args, like, like, like)
		argN += 3
	}

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM exercises WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.PageSize, offset)
	query := fmt.Sprintf(`
		SELECT id, trainer_id, name, muscle_group, equipment, video_url, notes, created_at, updated_at
		FROM exercises WHERE %s ORDER BY name ASC LIMIT $%d OFFSET $%d`,
		where, argN, argN+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var exercises []Exercise
	for rows.Next() {
		var e Exercise
		if err := rows.Scan(
			&e.ID, &e.TrainerID, &e.Name, &e.MuscleGroup, &e.Equipment, &e.VideoURL, &e.Notes, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		exercises = append(exercises, e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return exercises, total, nil
}

func (r *Repository) Update(ctx context.Context, e *Exercise) error {
	query := `
		UPDATE exercises SET name=$1, muscle_group=$2, equipment=$3, video_url=$4, notes=$5, updated_at=$6
		WHERE id=$7 AND trainer_id=$8`
	res, err := r.db.Exec(ctx, query,
		e.Name, e.MuscleGroup, e.Equipment, e.VideoURL, e.Notes, e.UpdatedAt, e.ID, e.TrainerID,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id, trainerID uuid.UUID) error {
	res, err := r.db.Exec(ctx, `DELETE FROM exercises WHERE id=$1 AND trainer_id=$2`, id, trainerID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
