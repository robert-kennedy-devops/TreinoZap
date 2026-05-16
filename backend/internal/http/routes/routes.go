package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/treinozap/backend/internal/admin"
	"github.com/treinozap/backend/internal/client"
	"github.com/treinozap/backend/internal/exercise"
	httpmiddleware "github.com/treinozap/backend/internal/http/middleware"
	"github.com/treinozap/backend/internal/http/response"
	"github.com/treinozap/backend/internal/message"
	"github.com/treinozap/backend/internal/trainer"
	"github.com/treinozap/backend/internal/workout"
)

type Handlers struct {
	Trainer  *trainer.Handler
	Client   *client.Handler
	Exercise *exercise.Handler
	Workout  *workout.Handler
	Message  *message.Handler
	Admin    *admin.Handler
}

func New(h Handlers, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://frontend:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public
		r.Post("/auth/register", h.Trainer.Register)
		r.Post("/auth/login", h.Trainer.Login)

		// Authenticated
		r.Group(func(r chi.Router) {
			r.Use(httpmiddleware.Auth(jwtSecret))

			r.Get("/me", h.Trainer.Me)

			r.Route("/clients", func(r chi.Router) {
				r.Get("/", h.Client.List)
				r.Post("/", h.Client.Create)
				r.Get("/{id}", h.Client.GetByID)
				r.Put("/{id}", h.Client.Update)
				r.Delete("/{id}", h.Client.Delete)

				r.Get("/{clientId}/workouts", h.Workout.ListByClient)
				r.Post("/{clientId}/workouts", h.Workout.Create)
				r.Get("/{clientId}/messages", h.Message.ListByClient)
			})

			r.Route("/exercises", func(r chi.Router) {
				r.Get("/", h.Exercise.List)
				r.Post("/", h.Exercise.Create)
				r.Get("/{id}", h.Exercise.GetByID)
				r.Put("/{id}", h.Exercise.Update)
				r.Delete("/{id}", h.Exercise.Delete)
			})

			r.Route("/workouts", func(r chi.Router) {
				r.Get("/{id}", h.Workout.GetByID)
				r.Put("/{id}", h.Workout.Update)
				r.Delete("/{id}", h.Workout.Archive)
				r.Post("/{id}/activate", h.Workout.Activate)
				r.Post("/{id}/send-whatsapp", h.Workout.SendWhatsApp)
			})

			r.Get("/messages", h.Message.List)

			// Admin only
			r.Route("/admin", func(r chi.Router) {
				r.Use(httpmiddleware.RequireAdmin())
				r.Route("/whatsapp", func(r chi.Router) {
					r.Get("/status", h.Admin.Status)
					r.Get("/qr", h.Admin.QRCode)
					r.Post("/connect", h.Admin.Connect)
					r.Post("/disconnect", h.Admin.Disconnect)
				})
				r.Get("/trainers", h.Admin.ListTrainers)
				r.Get("/clients", h.Admin.ListClients)
			})
		})
	})

	return r
}
