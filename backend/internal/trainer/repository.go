package trainer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("treinador não encontrado")
var ErrEmailInUse = errors.New("e-mail já cadastrado")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, t *Trainer, passwordHash string) error {
	query := `
		INSERT INTO trainers (id, name, email, password_hash, phone, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(ctx, query,
		t.ID, t.Name, t.Email, passwordHash, t.Phone, t.Role, t.Status, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*Trainer, string, error) {
	query := `
		SELECT id, name, email, password_hash, phone, role, status, created_at, updated_at
		FROM trainers WHERE email = $1`

	var t Trainer
	var hash string
	err := r.db.QueryRow(ctx, query, email).Scan(
		&t.ID, &t.Name, &t.Email, &hash, &t.Phone, &t.Role, &t.Status, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, "", ErrNotFound
	}
	if err != nil {
		return nil, "", err
	}
	return &t, hash, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Trainer, error) {
	query := `
		SELECT id, name, email, phone, role, status, created_at, updated_at
		FROM trainers WHERE id = $1`

	var t Trainer
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.Email, &t.Phone, &t.Role, &t.Status, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &t, err
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM trainers WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

func (r *Repository) ListAll(ctx context.Context) ([]Trainer, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, email, phone, role, status, created_at, updated_at
		FROM trainers ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainers []Trainer
	for rows.Next() {
		var t Trainer
		if err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Phone, &t.Role, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		trainers = append(trainers, t)
	}
	if trainers == nil {
		trainers = []Trainer{}
	}
	return trainers, rows.Err()
}
