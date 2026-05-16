package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/treinozap/backend/internal/admin"
	"github.com/treinozap/backend/internal/automation"
	"github.com/treinozap/backend/internal/client"
	"github.com/treinozap/backend/internal/config"
	"github.com/treinozap/backend/internal/database"
	"github.com/treinozap/backend/internal/exercise"
	"github.com/treinozap/backend/internal/http/routes"
	"github.com/treinozap/backend/internal/message"
	"github.com/treinozap/backend/internal/trainer"
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
	log.Println("conectado ao banco de dados")

	// Repositories
	trainerRepo := trainer.NewRepository(db)
	clientRepo := client.NewRepository(db)
	exerciseRepo := exercise.NewRepository(db)
	workoutRepo := workout.NewRepository(db)
	messageRepo := message.NewRepository(db)

	// Services
	trainerSvc := trainer.NewService(trainerRepo, cfg)
	clientSvc := client.NewService(clientRepo)
	exerciseSvc := exercise.NewService(exerciseRepo)
	workoutSvc := workout.NewService(workoutRepo)

	// WhatsApp
	var sender whatsapp.Sender
	var connMgr whatsapp.ConnectionManager

	switch cfg.WhatsAppProvider {
	case "whatsmeow":
		log.Println("WhatsApp provider: whatsmeow")
		var wa *whatsapp.WhatsMeowClient
		wa, err = whatsapp.NewWhatsMeowClient(cfg, func(fromPhone, text string) {
			// wa is always set by the time a message arrives (after Connect)
			h := automation.NewIncomingHandler(clientSvc, workoutSvc, messageRepo, wa)
			if err := h.Handle(context.Background(), fromPhone, text); err != nil {
				log.Printf("[whatsapp] erro ao processar mensagem de %s: %v", fromPhone, err)
			}
		})
		if err != nil {
			log.Fatalf("erro ao inicializar whatsmeow: %v", err)
		}
		sender = wa
		connMgr = wa
		// Auto-connect só se já houver sessão pareada salva no banco.
		// Sem sessão: aguarda o admin clicar "Conectar" no painel (/settings/whatsapp).
		if wa.HasSession() {
			if connectErr := wa.Connect(ctx); connectErr != nil {
				log.Printf("[whatsmeow] erro ao reconectar sessão existente: %v", connectErr)
			}
		} else {
			log.Println("[whatsmeow] sem sessão salva — acesse /settings/whatsapp para escanear o QR Code")
		}
	default:
		log.Println("WhatsApp provider: mock")
		sender = whatsapp.NewMockSender()
		connMgr = whatsapp.NewMockConnectionManager()
	}

	// Handlers
	h := routes.Handlers{
		Trainer:  trainer.NewHandler(trainerSvc),
		Client:   client.NewHandler(clientSvc),
		Exercise: exercise.NewHandler(exerciseSvc),
		Workout:  workout.NewHandler(workoutSvc, clientSvc, messageRepo, sender),
		Message:  message.NewHandler(messageRepo, clientSvc),
		Admin:    admin.NewHandler(connMgr, trainerSvc, clientSvc, cfg.WhatsAppAdminEnabled),
	}

	router := routes.New(h, cfg.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("backend iniciado em http://localhost%s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("servidor HTTP: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("encerrando servidor...")
	shutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Printf("erro ao encerrar: %v", err)
	}
	log.Println("servidor encerrado")
}
