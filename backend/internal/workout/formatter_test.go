package workout_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/treinozap/backend/internal/workout"
)

func restPtr(v int) *int { return &v }

func sampleWorkout() *workout.Workout {
	now := time.Now()
	sectionID := uuid.New()
	return &workout.Workout{
		ID:        uuid.New(),
		TrainerID: uuid.New(),
		ClientID:  uuid.New(),
		Name:      "Plano Hipertrofia",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
		Sections: []workout.Section{
			{
				ID:         sectionID,
				WorkoutID:  uuid.New(),
				Name:       "Treino A - Peito",
				OrderIndex: 1,
				CreatedAt:  now,
				UpdatedAt:  now,
				Exercises: []workout.WorkoutExercise{
					{
						ID:           uuid.New(),
						SectionID:    sectionID,
						ExerciseName: "Supino reto",
						Sets:         "4",
						Reps:         "10",
						RestSeconds:  restPtr(60),
						LoadNote:     "Carga moderada",
						OrderIndex:   1,
						CreatedAt:    now,
						UpdatedAt:    now,
					},
					{
						ID:           uuid.New(),
						SectionID:    sectionID,
						ExerciseName: "Supino inclinado",
						Sets:         "3",
						Reps:         "12",
						OrderIndex:   2,
						CreatedAt:    now,
						UpdatedAt:    now,
					},
				},
			},
		},
	}
}

func TestFormatWorkoutForWhatsApp(t *testing.T) {
	w := sampleWorkout()
	result := workout.FormatWorkoutForWhatsApp("João", w)

	if !strings.Contains(result, "João") {
		t.Error("deve conter o nome do cliente")
	}
	if !strings.Contains(result, "*Plano Hipertrofia*") {
		t.Error("deve conter o nome do treino em negrito")
	}
	if !strings.Contains(result, "Supino reto") {
		t.Error("deve conter o exercício")
	}
	if !strings.Contains(result, "4x10") {
		t.Error("deve conter séries x repetições")
	}
	if !strings.Contains(result, "descanso 60s") {
		t.Error("deve conter o descanso")
	}
	if !strings.Contains(result, "Carga moderada") {
		t.Error("deve conter a observação de carga")
	}
	if strings.Contains(result, "descanso 0s") {
		t.Error("não deve exibir descanso quando zero")
	}
}

func TestFormatWorkoutSectionForWhatsApp(t *testing.T) {
	w := sampleWorkout()
	result := workout.FormatWorkoutSectionForWhatsApp("Maria", w, "treino a")

	if !strings.Contains(result, "Maria") {
		t.Error("deve conter o nome do cliente")
	}
	if !strings.Contains(result, "Treino A - Peito") {
		t.Error("deve conter o nome da seção")
	}
	if !strings.Contains(result, "Supino reto") {
		t.Error("deve conter o exercício da seção")
	}
}

func TestFormatWorkoutSectionNotFound(t *testing.T) {
	w := sampleWorkout()
	result := workout.FormatWorkoutSectionForWhatsApp("Maria", w, "treino d")

	if !strings.Contains(result, "Não encontrei") {
		t.Error("deve informar que a seção não foi encontrada")
	}
}
