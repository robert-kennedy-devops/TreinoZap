package workout

import "github.com/google/uuid"

type CreateExerciseInput struct {
	ExerciseID    *uuid.UUID `json:"exercise_id"`
	ExerciseName  string     `json:"exercise_name"`
	Sets          string     `json:"sets"`
	Reps          string     `json:"reps"`
	RestSeconds   *int       `json:"rest_seconds"`
	LoadNote      string     `json:"load_note"`
	TechniqueNote string     `json:"technique_note"`
	VideoURL      string     `json:"video_url"`
	OrderIndex    int        `json:"order_index"`
}

type CreateSectionInput struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	OrderIndex  int                   `json:"order_index"`
	Exercises   []CreateExerciseInput `json:"exercises"`
}

type CreateRequest struct {
	Name     string               `json:"name"`
	StartsAt string               `json:"starts_at"`
	EndsAt   string               `json:"ends_at"`
	Sections []CreateSectionInput `json:"sections"`
}

type UpdateRequest struct {
	Name     string               `json:"name"`
	StartsAt string               `json:"starts_at"`
	EndsAt   string               `json:"ends_at"`
	Status   string               `json:"status"`
	Sections []CreateSectionInput `json:"sections"`
}
