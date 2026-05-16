package workout

import (
	"context"
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

func (s *Service) Create(ctx context.Context, trainerID, clientID uuid.UUID, req CreateRequest) (*Workout, error) {
	now := time.Now().UTC()
	w := &Workout{
		ID:        uuid.New(),
		TrainerID: trainerID,
		ClientID:  clientID,
		Name:      strings.TrimSpace(req.Name),
		Status:    "draft",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.StartsAt != "" {
		w.StartsAt = &req.StartsAt
	}
	if req.EndsAt != "" {
		w.EndsAt = &req.EndsAt
	}

	w.Sections = buildSections(req.Sections, now)
	return w, s.repo.Create(ctx, w)
}

func (s *Service) GetByID(ctx context.Context, id, trainerID uuid.UUID) (*Workout, error) {
	return s.repo.FindByID(ctx, id, trainerID)
}

func (s *Service) GetActiveByClientID(ctx context.Context, clientID uuid.UUID) (*Workout, error) {
	return s.repo.FindActiveByClientID(ctx, clientID)
}

func (s *Service) ListByClient(ctx context.Context, clientID, trainerID uuid.UUID) ([]Workout, error) {
	workouts, err := s.repo.ListByClient(ctx, clientID, trainerID)
	if workouts == nil {
		workouts = []Workout{}
	}
	return workouts, err
}

func (s *Service) Activate(ctx context.Context, id, trainerID uuid.UUID) (*Workout, error) {
	w, err := s.repo.FindByID(ctx, id, trainerID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Activate(ctx, id, w.ClientID, trainerID); err != nil {
		return nil, err
	}
	w.Status = "active"
	return w, nil
}

func (s *Service) Update(ctx context.Context, id, trainerID uuid.UUID, req UpdateRequest) (*Workout, error) {
	w, err := s.repo.FindByID(ctx, id, trainerID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	w.Name = strings.TrimSpace(req.Name)
	w.UpdatedAt = now

	if req.StartsAt != "" {
		w.StartsAt = &req.StartsAt
	}
	if req.EndsAt != "" {
		w.EndsAt = &req.EndsAt
	}
	if req.Status != "" && req.Status != "active" {
		w.Status = req.Status
	}

	w.Sections = buildSections(req.Sections, now)
	return w, s.repo.UpdateWithSections(ctx, w)
}

func (s *Service) Archive(ctx context.Context, id, trainerID uuid.UUID) error {
	return s.repo.Archive(ctx, id, trainerID)
}

func buildSections(inputs []CreateSectionInput, now time.Time) []Section {
	sections := make([]Section, 0, len(inputs))
	for _, si := range inputs {
		sectionID := uuid.New()
		exercises := make([]WorkoutExercise, 0, len(si.Exercises))
		for _, ei := range si.Exercises {
			name := strings.TrimSpace(ei.ExerciseName)
			if name == "" {
				name = "Exercício"
			}
			exercises = append(exercises, WorkoutExercise{
				ID:            uuid.New(),
				SectionID:     sectionID,
				ExerciseID:    ei.ExerciseID,
				ExerciseName:  name,
				Sets:          ei.Sets,
				Reps:          ei.Reps,
				RestSeconds:   ei.RestSeconds,
				LoadNote:      ei.LoadNote,
				TechniqueNote: ei.TechniqueNote,
				VideoURL:      ei.VideoURL,
				OrderIndex:    ei.OrderIndex,
				CreatedAt:     now,
				UpdatedAt:     now,
			})
		}
		sections = append(sections, Section{
			ID:          sectionID,
			Name:        strings.TrimSpace(si.Name),
			Description: strings.TrimSpace(si.Description),
			OrderIndex:  si.OrderIndex,
			Exercises:   exercises,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return sections
}
