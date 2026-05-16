package whatsapp

import (
	"context"
	"errors"
	"log"
	"time"
)

// MockSender logs messages instead of sending them via WhatsApp.
type MockSender struct{}

func NewMockSender() *MockSender {
	return &MockSender{}
}

func (s *MockSender) SendText(_ context.Context, phone, message string) error {
	log.Printf("[WhatsApp MOCK] → %s: %s", phone, message)
	return nil
}

// MockConnectionManager simulates a WhatsApp connection.
type MockConnectionManager struct {
	connected       bool
	lastConnectedAt time.Time
}

func NewMockConnectionManager() *MockConnectionManager {
	return &MockConnectionManager{}
}

func (m *MockConnectionManager) Connect(_ context.Context) error {
	m.connected = true
	m.lastConnectedAt = time.Now().UTC()
	log.Println("[WhatsApp MOCK] conectado")
	return nil
}

func (m *MockConnectionManager) Disconnect(_ context.Context) error {
	m.connected = false
	log.Println("[WhatsApp MOCK] desconectado")
	return nil
}

func (m *MockConnectionManager) Status(_ context.Context) (Status, error) {
	if m.connected {
		st := Status{Connected: true, Phone: "mock", JID: "mock@s.whatsapp.net"}
		if !m.lastConnectedAt.IsZero() {
			st.LastConnected = m.lastConnectedAt.Format(time.RFC3339)
		}
		return st, nil
	}
	st := Status{Connected: false}
	if !m.lastConnectedAt.IsZero() {
		st.LastConnected = m.lastConnectedAt.Format(time.RFC3339)
	}
	return st, nil
}

func (m *MockConnectionManager) QRCode(_ context.Context) (string, error) {
	if m.connected {
		return "", errors.New("já conectado")
	}
	return "MOCK-QR-CODE-DATA", nil
}
