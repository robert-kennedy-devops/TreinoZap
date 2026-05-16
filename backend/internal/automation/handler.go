package automation

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/google/uuid"
	"github.com/treinozap/backend/internal/client"
	"github.com/treinozap/backend/internal/message"
	"github.com/treinozap/backend/internal/whatsapp"
	"github.com/treinozap/backend/internal/workout"
)

type IncomingHandler struct {
	clientSvc   *client.Service
	workoutSvc  *workout.Service
	messageRepo *message.Repository
	sender      whatsapp.Sender
}

func NewIncomingHandler(
	clientSvc *client.Service,
	workoutSvc *workout.Service,
	messageRepo *message.Repository,
	sender whatsapp.Sender,
) *IncomingHandler {
	return &IncomingHandler{
		clientSvc:   clientSvc,
		workoutSvc:  workoutSvc,
		messageRepo: messageRepo,
		sender:      sender,
	}
}

func (h *IncomingHandler) Handle(ctx context.Context, fromPhone, text string) error {
	phone := normalizePhone(fromPhone)
	cmd := normalizeCommand(text)

	// Save inbound message before any processing
	inbound := message.NewInbound(nil, nil, phone, text, cmd)

	c, err := h.clientSvc.FindByPhone(ctx, phone)
	if err != nil {
		// Save inbound without client link
		_ = h.messageRepo.Save(ctx, inbound)
		return h.send(ctx, phone, nil, nil,
			"Não encontrei seu cadastro. Fale com seu treinador para liberar seu acesso.",
			"unknown_client",
		)
	}

	// Update inbound message with client/trainer info
	inbound.ClientID = &c.ID
	inbound.TrainerID = &c.TrainerID
	_ = h.messageRepo.Save(ctx, inbound)

	if c.Status == "inactive" {
		return h.send(ctx, phone, &c.TrainerID, &c.ID,
			"Seu acesso está inativo. Fale com seu treinador.",
			"inactive_client",
		)
	}
	if c.Status == "blocked" {
		log.Printf("[automation] cliente bloqueado ignorado: %s", phone)
		return nil
	}

	trainerID := c.TrainerID
	clientID := c.ID

	switch cmd {
	case "oi", "ola", "ola!", "oi!", "bom dia", "boa tarde", "boa noite", "hello", "hi":
		return h.sendMenu(ctx, &trainerID, &clientID, phone, c.Name)

	case "treino", "meu treino", "quero meu treino", "enviar treino":
		return h.sendFullWorkout(ctx, &trainerID, &clientID, phone, c.Name)

	case "status", "tenho treino", "tem treino", "meu status":
		return h.sendStatus(ctx, &trainerID, &clientID, phone, c.Name)

	case "menu":
		return h.sendMenu(ctx, &trainerID, &clientID, phone, c.Name)

	case "a", "treino a":
		return h.sendSection(ctx, &trainerID, &clientID, phone, c.Name, "treino a")

	case "b", "treino b":
		return h.sendSection(ctx, &trainerID, &clientID, phone, c.Name, "treino b")

	case "c", "treino c":
		return h.sendSection(ctx, &trainerID, &clientID, phone, c.Name, "treino c")

	case "ajuda":
		return h.sendHelp(ctx, &trainerID, &clientID, phone)

	default:
		return h.send(ctx, phone, &trainerID, &clientID,
			fmt.Sprintf("Olá, %s! Não entendi sua mensagem. Envie *menu* para ver as opções disponíveis.", c.Name),
			"unknown_command",
		)
	}
}

func (h *IncomingHandler) sendFullWorkout(ctx context.Context, trainerID, clientID *uuid.UUID, phone, clientName string) error {
	w, err := h.workoutSvc.GetActiveByClientID(ctx, *clientID)
	if err != nil {
		return h.send(ctx, phone, trainerID, clientID,
			"Você ainda não tem um treino ativo. Fale com seu treinador.",
			"no_active_workout",
		)
	}
	text := workout.FormatWorkoutForWhatsApp(clientName, w)
	return h.send(ctx, phone, trainerID, clientID, text, "treino")
}

func (h *IncomingHandler) sendSection(ctx context.Context, trainerID, clientID *uuid.UUID, phone, clientName, section string) error {
	w, err := h.workoutSvc.GetActiveByClientID(ctx, *clientID)
	if err != nil {
		return h.send(ctx, phone, trainerID, clientID,
			"Você ainda não tem um treino ativo. Fale com seu treinador.",
			"no_active_workout",
		)
	}
	text := workout.FormatWorkoutSectionForWhatsApp(clientName, w, section)
	return h.send(ctx, phone, trainerID, clientID, text, section)
}

func (h *IncomingHandler) sendStatus(ctx context.Context, trainerID, clientID *uuid.UUID, phone, clientName string) error {
	w, err := h.workoutSvc.GetActiveByClientID(ctx, *clientID)
	if err != nil {
		return h.send(ctx, phone, trainerID, clientID,
			fmt.Sprintf("Olá, %s! Você ainda não tem um treino ativo no momento. Aguarde seu treinador cadastrar seu treino.", clientName),
			"status_no_workout",
		)
	}
	text := fmt.Sprintf(
		"✅ Olá, %s! Você tem um treino ativo:\n\n*%s*\n\nEnvie *treino* para receber o treino completo ou *menu* para ver todas as opções.",
		clientName, w.Name,
	)
	return h.send(ctx, phone, trainerID, clientID, text, "status")
}

func (h *IncomingHandler) sendMenu(ctx context.Context, trainerID, clientID *uuid.UUID, phone, clientName string) error {
	text := fmt.Sprintf(
		"Olá, %s! 👋\n\nEnvie uma das opções:\n\n*status* - verificar se tem treino ativo\n*treino* - receber treino completo\n*A* - receber Treino A\n*B* - receber Treino B\n*C* - receber Treino C\n*ajuda* - ver instruções",
		clientName,
	)
	return h.send(ctx, phone, trainerID, clientID, text, "menu")
}

func (h *IncomingHandler) sendHelp(ctx context.Context, trainerID, clientID *uuid.UUID, phone string) error {
	text := "Comandos disponíveis:\n\n*status* - verificar se tem treino ativo\n*treino* - receber treino completo\n*A*, *B*, *C* - receber seção específica\n*menu* - ver opções\n\nEm caso de dúvidas, fale com seu treinador."
	return h.send(ctx, phone, trainerID, clientID, text, "ajuda")
}

func (h *IncomingHandler) send(ctx context.Context, phone string, trainerID, clientID *uuid.UUID, text, cmd string) error {
	if err := h.sender.SendText(ctx, phone, text); err != nil {
		log.Printf("[automation] erro ao enviar para %s: %v", phone, err)
		return err
	}
	out := message.NewOutbound(trainerID, clientID, phone, text, cmd)
	return h.messageRepo.Save(ctx, out)
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

func normalizeCommand(text string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	clean, _, _ := transform.String(t, text)
	return strings.ToLower(strings.TrimSpace(clean))
}
