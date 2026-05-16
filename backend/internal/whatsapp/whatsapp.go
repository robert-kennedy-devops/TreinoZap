package whatsapp

import "context"

type Status struct {
	Connected     bool   `json:"connected"`
	Phone         string `json:"phone,omitempty"`
	JID           string `json:"jid,omitempty"`
	LastConnected string `json:"last_connected,omitempty"`
}

// Sender sends WhatsApp messages.
type Sender interface {
	SendText(ctx context.Context, phone string, message string) error
}

// ConnectionManager manages the WhatsApp connection lifecycle.
type ConnectionManager interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Status(ctx context.Context) (Status, error)
	QRCode(ctx context.Context) (string, error)
}
