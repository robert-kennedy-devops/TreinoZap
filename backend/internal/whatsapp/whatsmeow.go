package whatsapp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	whatsmeow "go.mau.fi/whatsmeow"

	"github.com/treinozap/backend/internal/config"
)

// OnMessageFunc is called when a private text message is received.
type OnMessageFunc func(fromPhone, text string)

// WhatsMeowClient wraps whatsmeow and implements Sender and ConnectionManager.
type WhatsMeowClient struct {
	mu              sync.RWMutex
	client          *whatsmeow.Client
	container       *sqlstore.Container
	onMessage       OnMessageFunc
	qrCode          string
	lastConnectedAt time.Time
}

func NewWhatsMeowClient(cfg *config.Config, onMessage OnMessageFunc) (*WhatsMeowClient, error) {
	logger := waLog.Stdout("whatsmeow", "WARN", true)

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("abrir db para whatsmeow: %w", err)
	}

	container := sqlstore.NewWithDB(db, "postgres", logger)
	if err := container.Upgrade(context.Background()); err != nil {
		return nil, fmt.Errorf("upgrade whatsmeow store: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("obter device: %w", err)
	}

	client := whatsmeow.NewClient(deviceStore, logger)

	c := &WhatsMeowClient{
		client:    client,
		container: container,
		onMessage: onMessage,
	}

	client.AddEventHandler(c.handleEvent)
	return c, nil
}

func (c *WhatsMeowClient) handleEvent(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		c.handleMessage(v)

	case *events.Connected:
		log.Println("[whatsmeow] conectado")
		c.mu.Lock()
		c.qrCode = ""
		c.lastConnectedAt = time.Now().UTC()
		c.mu.Unlock()

	case *events.Disconnected:
		log.Println("[whatsmeow] desconectado — tentando reconectar em 10s")
		go func() {
			time.Sleep(10 * time.Second)
			c.mu.RLock()
			cl := c.client
			c.mu.RUnlock()
			// Only reconnect if we have a stored session (device was paired)
			if cl.Store.ID == nil {
				return
			}
			log.Println("[whatsmeow] reconectando...")
			if err := cl.Connect(); err != nil {
				log.Printf("[whatsmeow] erro ao reconectar: %v", err)
			}
		}()

	case *events.LoggedOut:
		log.Println("[whatsmeow] sessão encerrada pelo servidor — escaneie o QR novamente")
		c.mu.Lock()
		c.qrCode = ""
		c.mu.Unlock()
	}
}

func (c *WhatsMeowClient) handleMessage(evt *events.Message) {
	info := evt.Info

	if info.IsGroup || info.IsFromMe || info.Chat.Server == types.BroadcastServer {
		return
	}

	text := extractText(evt.Message)
	if strings.TrimSpace(text) == "" {
		return
	}

	// When the sender JID uses the LID server ("lid"), AddressingMode may be
	// empty (whatsmeow bug). Detect by server name and use SenderAlt for the
	// real phone-number JID, which is what the client DB stores.
	phone := info.Sender.User
	if info.Sender.Server == types.HiddenUserServer && !info.SenderAlt.IsEmpty() {
		phone = info.SenderAlt.User
	}
	log.Printf("[whatsmeow] mensagem de %s: %s", phone, text)

	if c.onMessage != nil {
		c.onMessage(phone, text)
	}
}

func extractText(msg *waProto.Message) string {
	if msg == nil {
		return ""
	}
	if msg.Conversation != nil {
		return *msg.Conversation
	}
	if msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.Text != nil {
		return *msg.ExtendedTextMessage.Text
	}
	return ""
}

// SendText implements Sender.
func (c *WhatsMeowClient) SendText(ctx context.Context, phone, message string) error {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	if client == nil || !client.IsConnected() {
		return errors.New("whatsapp não está conectado")
	}

	jid := types.NewJID(phone, types.DefaultUserServer)
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err := client.SendMessage(ctx, jid, msg)
	return err
}

// HasSession returns true if a paired device session exists in the store.
func (c *WhatsMeowClient) HasSession() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client.Store.ID != nil
}

// Connect implements ConnectionManager.
func (c *WhatsMeowClient) Connect(ctx context.Context) error {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	if client.Store.ID == nil {
		// No session — need the QR pairing flow.
		// GetQRChannel MUST be called before client.Connect(), so if the
		// WebSocket is already open (e.g. from a previous attempt) we must
		// disconnect cleanly first.
		if client.IsConnected() {
			client.Disconnect()
			time.Sleep(300 * time.Millisecond)
		}

		// Use context.Background() so the QR session is NOT tied to the HTTP
		// request context. When the handler responds, r.Context() is cancelled
		// and would kill the QR channel before the user scans it.
		qrChan, err := client.GetQRChannel(context.Background())
		if err != nil {
			return fmt.Errorf("obter canal QR: %w", err)
		}

		go func() {
			if err := client.Connect(); err != nil {
				log.Printf("[whatsmeow] erro ao conectar: %v", err)
			}
		}()

		go func() {
			for item := range qrChan {
				if item.Event == whatsmeow.QRChannelEventCode {
					c.mu.Lock()
					c.qrCode = item.Code
					c.mu.Unlock()
					log.Println("[whatsmeow] novo QR gerado — escaneie no painel admin")
				} else if item == whatsmeow.QRChannelSuccess {
					log.Println("[whatsmeow] pareamento realizado com sucesso")
				} else {
					log.Printf("[whatsmeow] QR evento: %s", item.Event)
				}
			}
		}()

		return nil
	}

	// Session exists — reconnect directly (no QR needed).
	if client.IsConnected() {
		return nil
	}
	return client.Connect()
}

// Disconnect implements ConnectionManager.
func (c *WhatsMeowClient) Disconnect(_ context.Context) error {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	client.Disconnect()
	log.Println("[whatsmeow] desconectado pelo admin")
	return nil
}

// authenticated returns true only when the WebSocket is open AND the device
// is fully paired (Store.ID is set). A bare WebSocket during the QR flow is
// NOT considered authenticated.
func (c *WhatsMeowClient) authenticated() bool {
	return c.client.IsConnected() && c.client.Store.ID != nil
}

// Status implements ConnectionManager.
func (c *WhatsMeowClient) Status(_ context.Context) (Status, error) {
	c.mu.RLock()
	client := c.client
	lastConnectedAt := c.lastConnectedAt
	c.mu.RUnlock()

	auth := client.IsConnected() && client.Store.ID != nil
	st := Status{Connected: auth}
	if client.Store.ID != nil {
		st.JID = client.Store.ID.String()
		st.Phone = client.Store.ID.User
	}
	if !lastConnectedAt.IsZero() {
		st.LastConnected = lastConnectedAt.Format(time.RFC3339)
	}
	return st, nil
}

// QRCode implements ConnectionManager.
func (c *WhatsMeowClient) QRCode(_ context.Context) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Fully authenticated — no QR needed
	if c.client.IsConnected() && c.client.Store.ID != nil {
		return "", errors.New("já conectado")
	}
	if c.qrCode == "" {
		return "", errors.New("QR Code não disponível — clique em Conectar primeiro")
	}
	return c.qrCode, nil
}
