package message

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID                uuid.UUID  `json:"id"`
	ChannelID         *uuid.UUID `json:"channel_id,omitempty"`
	TrainerID         *uuid.UUID `json:"trainer_id,omitempty"`
	ClientID          *uuid.UUID `json:"client_id,omitempty"`
	Direction         string     `json:"direction"`
	Phone             string     `json:"phone"`
	Message           string     `json:"message"`
	Command           string     `json:"command,omitempty"`
	Status            string     `json:"status,omitempty"`
	ProviderMessageID string     `json:"provider_message_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

func NewOutbound(trainerID, clientID *uuid.UUID, phone, text, command string) *Message {
	return &Message{
		ID:        uuid.New(),
		TrainerID: trainerID,
		ClientID:  clientID,
		Direction: "outbound",
		Phone:     phone,
		Message:   text,
		Command:   command,
		Status:    "sent",
		CreatedAt: time.Now().UTC(),
	}
}

func NewInbound(trainerID, clientID *uuid.UUID, phone, text, command string) *Message {
	return &Message{
		ID:        uuid.New(),
		TrainerID: trainerID,
		ClientID:  clientID,
		Direction: "inbound",
		Phone:     phone,
		Message:   text,
		Command:   command,
		Status:    "received",
		CreatedAt: time.Now().UTC(),
	}
}
