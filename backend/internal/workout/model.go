package workout

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	ID        uuid.UUID `json:"id"`
	TrainerID uuid.UUID `json:"trainer_id"`
	ClientID  uuid.UUID `json:"client_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartsAt  *string   `json:"starts_at,omitempty"`
	EndsAt    *string   `json:"ends_at,omitempty"`
	Sections  []Section `json:"sections,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Section struct {
	ID          uuid.UUID         `json:"id"`
	WorkoutID   uuid.UUID         `json:"workout_id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	OrderIndex  int               `json:"order_index"`
	Exercises   []WorkoutExercise `json:"exercises,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type WorkoutExercise struct {
	ID            uuid.UUID  `json:"id"`
	SectionID     uuid.UUID  `json:"section_id"`
	ExerciseID    *uuid.UUID `json:"exercise_id,omitempty"`
	ExerciseName  string     `json:"exercise_name"`
	Sets          string     `json:"sets,omitempty"`
	Reps          string     `json:"reps,omitempty"`
	RestSeconds   *int       `json:"rest_seconds,omitempty"`
	LoadNote      string     `json:"load_note,omitempty"`
	TechniqueNote string     `json:"technique_note,omitempty"`
	VideoURL      string     `json:"video_url,omitempty"`
	OrderIndex    int        `json:"order_index"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
