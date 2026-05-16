package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/treinozap/backend/internal/automation"
	"github.com/treinozap/backend/internal/client"
	"github.com/treinozap/backend/internal/config"
	"github.com/treinozap/backend/internal/database"
	"github.com/treinozap/backend/internal/message"
	"github.com/treinozap/backend/internal/whatsapp"
	"github.com/treinozap/backend/internal/workout"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("conexão com banco: %v", err)
	}
	defer db.Close()

	clientRepo := client.NewRepository(db)
	clientSvc := client.NewService(clientRepo)
	workoutRepo := workout.NewRepository(db)
	messageRepo := message.NewRepository(db)
	workoutSvc := workout.NewService(workoutRepo)

	var sender whatsapp.Sender

	switch cfg.WhatsAppProvider {
	case "whatsmeow":
		log.Println("Worker: inicializando whatsmeow...")
		wa, err := whatsapp.NewWhatsMeowClient(cfg, func(fromPhone, text string) {
			handler := automation.NewIncomingHandler(clientSvc, workoutSvc, messageRepo, sender)
			if err := handler.Handle(context.Background(), fromPhone, text); err != nil {
				log.Printf("[worker] erro ao processar mensagem: %v", err)
			}
		})
		if err != nil {
			log.Fatalf("erro ao criar cliente whatsmeow: %v", err)
		}
		sender = wa
		if err := wa.Connect(ctx); err != nil {
			log.Printf("[worker] erro ao conectar: %v", err)
		}
	default:
		log.Println("Worker: usando mock sender")
		sender = whatsapp.NewMockSender()
	}

	_ = sender

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Worker iniciado. Aguardando mensagens...")
	<-quit
	log.Println("Worker encerrado.")
}
