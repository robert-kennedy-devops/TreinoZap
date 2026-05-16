package message

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, m *Message) error {
	query := `
		INSERT INTO whatsapp_messages
		(id, channel_id, trainer_id, client_id, direction, phone, message, command, status, provider_message_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.db.Exec(ctx, query,
		m.ID, m.ChannelID, m.TrainerID, m.ClientID, m.Direction,
		m.Phone, m.Message, m.Command, m.Status, m.ProviderMessageID, m.CreatedAt,
	)
	return err
}

type ListFilter struct {
	TrainerID *uuid.UUID
	ClientID  *uuid.UUID
	Direction string
	Page      int
	PageSize  int
}

func (r *Repository) List(ctx context.Context, f ListFilter) ([]Message, int, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 50
	}
	offset := (f.Page - 1) * f.PageSize

	where := "1=1"
	args := []any{}
	argN := 1

	if f.TrainerID != nil {
		where += fmt.Sprintf(" AND trainer_id=$%d", argN)
		args = append(args, *f.TrainerID)
		argN++
	}
	if f.ClientID != nil {
		where += fmt.Sprintf(" AND client_id=$%d", argN)
		args = append(args, *f.ClientID)
		argN++
	}
	if f.Direction != "" {
		where += fmt.Sprintf(" AND direction=$%d", argN)
		args = append(args, f.Direction)
		argN++
	}

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM whatsapp_messages WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.PageSize, offset)
	query := fmt.Sprintf(`
		SELECT id, channel_id, trainer_id, client_id, direction, phone, message, command, status, provider_message_id, created_at
		FROM whatsapp_messages WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, argN, argN+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(
			&m.ID, &m.ChannelID, &m.TrainerID, &m.ClientID, &m.Direction,
			&m.Phone, &m.Message, &m.Command, &m.Status, &m.ProviderMessageID, &m.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return messages, total, nil
}
