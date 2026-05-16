package workout

import (
	"fmt"
	"strings"
)

// FormatWorkoutForWhatsApp formats a complete workout for WhatsApp.
func FormatWorkoutForWhatsApp(clientName string, w *Workout) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Olá, %s! Aqui está seu treino ativo:\n\n", clientName)
	fmt.Fprintf(&sb, "*%s*\n", w.Name)

	for _, section := range w.Sections {
		sb.WriteString("\n")
		fmt.Fprintf(&sb, "*%s*\n", section.Name)
		for i, e := range section.Exercises {
			line := formatExerciseLine(i+1, e)
			sb.WriteString(line)
		}
	}

	sb.WriteString("\nBom treino! 💪")
	return sb.String()
}

// FormatWorkoutSectionForWhatsApp formats a single section for WhatsApp.
func FormatWorkoutSectionForWhatsApp(clientName string, w *Workout, sectionName string) string {
	var target *Section
	normalizedTarget := strings.ToLower(strings.TrimSpace(sectionName))

	for i := range w.Sections {
		if strings.Contains(strings.ToLower(w.Sections[i].Name), normalizedTarget) {
			target = &w.Sections[i]
			break
		}
	}

	if target == nil {
		return fmt.Sprintf("Olá, %s! Não encontrei a seção \"%s\" no seu treino ativo.", clientName, sectionName)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Olá, %s! Aqui está o *%s*:\n\n", clientName, target.Name)

	if target.Description != "" {
		fmt.Fprintf(&sb, "_%s_\n\n", target.Description)
	}

	for i, e := range target.Exercises {
		sb.WriteString(formatExerciseLine(i+1, e))
	}

	sb.WriteString("\nBom treino! 💪")
	return sb.String()
}

func formatExerciseLine(n int, e WorkoutExercise) string {
	var sb strings.Builder

	sets := e.Sets
	reps := e.Reps
	setsReps := ""
	if sets != "" && reps != "" {
		setsReps = fmt.Sprintf(" — %sx%s", sets, reps)
	} else if sets != "" {
		setsReps = fmt.Sprintf(" — %s séries", sets)
	} else if reps != "" {
		setsReps = fmt.Sprintf(" — %s reps", reps)
	}

	rest := ""
	if e.RestSeconds != nil && *e.RestSeconds > 0 {
		rest = fmt.Sprintf(" — descanso %ds", *e.RestSeconds)
	}

	fmt.Fprintf(&sb, "%d. %s%s%s\n", n, e.ExerciseName, setsReps, rest)

	if e.LoadNote != "" {
		fmt.Fprintf(&sb, "   Carga: %s\n", e.LoadNote)
	}
	if e.TechniqueNote != "" {
		fmt.Fprintf(&sb, "   Técnica: %s\n", e.TechniqueNote)
	}
	if e.VideoURL != "" {
		fmt.Fprintf(&sb, "   Vídeo: %s\n", e.VideoURL)
	}

	return sb.String()
}
