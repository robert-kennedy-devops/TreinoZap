package exercise

import (
	"time"

	"github.com/google/uuid"
)

type Exercise struct {
	ID          uuid.UUID `json:"id"`
	TrainerID   uuid.UUID `json:"trainer_id"`
	Name        string    `json:"name"`
	MuscleGroup string    `json:"muscle_group,omitempty"`
	Equipment   string    `json:"equipment,omitempty"`
	VideoURL    string    `json:"video_url,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
