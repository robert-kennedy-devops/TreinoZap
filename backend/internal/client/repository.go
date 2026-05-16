package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("cliente não encontrado")
var ErrPhoneInUse = errors.New("telefone já cadastrado")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, c *Client) error {
	query := `
		INSERT INTO clients (id, trainer_id, name, phone, status, goal, notes, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.Exec(ctx, query,
		c.ID, c.TrainerID, c.Name, c.Phone, c.Status, c.Goal, c.Notes, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrPhoneInUse
		}
	}
	return err
}

func (r *Repository) FindByID(ctx context.Context, id, trainerID uuid.UUID) (*Client, error) {
	query := `
		SELECT id, trainer_id, name, phone, status, goal, notes, created_at, updated_at
		FROM clients WHERE id=$1 AND trainer_id=$2`
	var c Client
	err := r.db.QueryRow(ctx, query, id, trainerID).Scan(
		&c.ID, &c.TrainerID, &c.Name, &c.Phone, &c.Status, &c.Goal, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &c, err
}

func (r *Repository) FindByPhone(ctx context.Context, phone string) (*Client, error) {
	query := `
		SELECT id, trainer_id, name, phone, status, goal, notes, created_at, updated_at
		FROM clients WHERE phone=$1`
	var c Client
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&c.ID, &c.TrainerID, &c.Name, &c.Phone, &c.Status, &c.Goal, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &c, err
}

type ListFilter struct {
	TrainerID uuid.UUID
	Search    string
	Page      int
	PageSize  int
}

func (r *Repository) List(ctx context.Context, f ListFilter) ([]Client, int, error) {
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
		where += fmt.Sprintf(" AND (name ILIKE $%d OR phone ILIKE $%d)", argN, argN+1)
		like := "%" + f.Search + "%"
		args = append(args, like, like)
		argN += 2
	}

	var total int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM clients WHERE "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, f.PageSize, offset)
	query := fmt.Sprintf(`
		SELECT id, trainer_id, name, phone, status, goal, notes, created_at, updated_at
		FROM clients WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, argN, argN+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var c Client
		if err := rows.Scan(
			&c.ID, &c.TrainerID, &c.Name, &c.Phone, &c.Status, &c.Goal, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		clients = append(clients, c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return clients, total, nil
}

func (r *Repository) Update(ctx context.Context, c *Client) error {
	query := `
		UPDATE clients SET name=$1, phone=$2, status=$3, goal=$4, notes=$5, updated_at=$6
		WHERE id=$7 AND trainer_id=$8`
	res, err := r.db.Exec(ctx, query,
		c.Name, c.Phone, c.Status, c.Goal, c.Notes, c.UpdatedAt, c.ID, c.TrainerID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrPhoneInUse
		}
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ListAllGlobal(ctx context.Context, search string) ([]Client, error) {
	where := "1=1"
	args := []any{}
	argN := 1

	if search != "" {
		where = fmt.Sprintf("(name ILIKE $%d OR phone ILIKE $%d)", argN, argN+1)
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT id, trainer_id, name, phone, status, goal, notes, created_at, updated_at
		FROM clients WHERE %s ORDER BY created_at DESC`, where), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.TrainerID, &c.Name, &c.Phone, &c.Status, &c.Goal, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	if clients == nil {
		clients = []Client{}
	}
	return clients, rows.Err()
}

func (r *Repository) SoftDelete(ctx context.Context, id, trainerID uuid.UUID) error {
	query := `UPDATE clients SET status='inactive', updated_at=NOW() WHERE id=$1 AND trainer_id=$2`
	res, err := r.db.Exec(ctx, query, id, trainerID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) PhoneExists(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	var exists bool
	var err error
	if excludeID != nil {
		err = r.db.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM clients WHERE phone=$1 AND id!=$2)`,
			phone, *excludeID,
		).Scan(&exists)
	} else {
		err = r.db.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM clients WHERE phone=$1)`,
			phone,
		).Scan(&exists)
	}
	return exists, err
}
