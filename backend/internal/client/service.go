package client

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, trainerID uuid.UUID, req CreateRequest) (*Client, error) {
	phone := normalizePhone(req.Phone)
	if phone == "" {
		return nil, errors.New("telefone inválido")
	}

	exists, err := s.repo.PhoneExists(ctx, phone, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPhoneInUse
	}

	now := time.Now().UTC()
	c := &Client{
		ID:        uuid.New(),
		TrainerID: trainerID,
		Name:      strings.TrimSpace(req.Name),
		Phone:     phone,
		Status:    "active",
		Goal:      strings.TrimSpace(req.Goal),
		Notes:     strings.TrimSpace(req.Notes),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetByID(ctx context.Context, id, trainerID uuid.UUID) (*Client, error) {
	return s.repo.FindByID(ctx, id, trainerID)
}

func (s *Service) List(ctx context.Context, f ListFilter) ([]Client, int, error) {
	clients, total, err := s.repo.List(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	if clients == nil {
		clients = []Client{}
	}
	return clients, total, nil
}

func (s *Service) Update(ctx context.Context, id, trainerID uuid.UUID, req UpdateRequest) (*Client, error) {
	c, err := s.repo.FindByID(ctx, id, trainerID)
	if err != nil {
		return nil, err
	}

	phone := normalizePhone(req.Phone)
	if phone == "" {
		return nil, errors.New("telefone inválido")
	}

	if phone != c.Phone {
		exists, err := s.repo.PhoneExists(ctx, phone, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrPhoneInUse
		}
	}

	c.Name = strings.TrimSpace(req.Name)
	c.Phone = phone
	c.Goal = strings.TrimSpace(req.Goal)
	c.Notes = strings.TrimSpace(req.Notes)
	c.UpdatedAt = time.Now().UTC()

	if req.Status != "" {
		if !isValidClientStatus(req.Status) {
			return nil, errors.New("status inválido")
		}
		c.Status = req.Status
	}

	return c, s.repo.Update(ctx, c)
}

func (s *Service) Delete(ctx context.Context, id, trainerID uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id, trainerID)
}

func (s *Service) ListAllGlobal(ctx context.Context, search string) ([]Client, error) {
	return s.repo.ListAllGlobal(ctx, search)
}

func (s *Service) FindByPhone(ctx context.Context, phone string) (*Client, error) {
	return s.repo.FindByPhone(ctx, normalizePhone(phone))
}

func normalizePhone(phone string) string {
	var b strings.Builder
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isValidClientStatus(status string) bool {
	switch status {
	case "active", "inactive", "blocked":
		return true
	default:
		return false
	}
}
