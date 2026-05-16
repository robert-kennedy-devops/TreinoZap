package trainer

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/treinozap/backend/internal/auth"
	"github.com/treinozap/backend/internal/config"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*Trainer, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailInUse
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := "trainer"
	if s.cfg.AdminEmail != "" && email == strings.ToLower(s.cfg.AdminEmail) {
		role = "admin"
	}

	now := time.Now().UTC()
	t := &Trainer{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(req.Name),
		Email:     email,
		Phone:     strings.TrimSpace(req.Phone),
		Role:      role,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, t, string(hash)); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	t, hash, err := s.repo.FindByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		return nil, errors.New("credenciais inválidas")
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	token, err := auth.GenerateToken(t.ID, t.Role, s.cfg.JWTSecret, s.cfg.JWTExpiresIn)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token, Trainer: *t}, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Trainer, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) ListAll(ctx context.Context) ([]Trainer, error) {
	return s.repo.ListAll(ctx)
}
