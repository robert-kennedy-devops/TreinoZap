package exercise

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

func (s *Service) Create(ctx context.Context, trainerID uuid.UUID, req CreateRequest) (*Exercise, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("nome é obrigatório")
	}
	now := time.Now().UTC()
	e := &Exercise{
		ID:          uuid.New(),
		TrainerID:   trainerID,
		Name:        strings.TrimSpace(req.Name),
		MuscleGroup: strings.TrimSpace(req.MuscleGroup),
		Equipment:   strings.TrimSpace(req.Equipment),
		VideoURL:    strings.TrimSpace(req.VideoURL),
		Notes:       strings.TrimSpace(req.Notes),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return e, s.repo.Create(ctx, e)
}

func (s *Service) GetByID(ctx context.Context, id, trainerID uuid.UUID) (*Exercise, error) {
	return s.repo.FindByID(ctx, id, trainerID)
}

func (s *Service) List(ctx context.Context, f ListFilter) ([]Exercise, int, error) {
	exercises, total, err := s.repo.List(ctx, f)
	if exercises == nil {
		exercises = []Exercise{}
	}
	return exercises, total, err
}

func (s *Service) Update(ctx context.Context, id, trainerID uuid.UUID, req UpdateRequest) (*Exercise, error) {
	e, err := s.repo.FindByID(ctx, id, trainerID)
	if err != nil {
		return nil, err
	}
	e.Name = strings.TrimSpace(req.Name)
	e.MuscleGroup = strings.TrimSpace(req.MuscleGroup)
	e.Equipment = strings.TrimSpace(req.Equipment)
	e.VideoURL = strings.TrimSpace(req.VideoURL)
	e.Notes = strings.TrimSpace(req.Notes)
	e.UpdatedAt = time.Now().UTC()
	return e, s.repo.Update(ctx, e)
}

func (s *Service) Delete(ctx context.Context, id, trainerID uuid.UUID) error {
	return s.repo.Delete(ctx, id, trainerID)
}
