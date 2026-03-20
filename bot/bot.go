package bot

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"whatsapp-bot/config"
	"whatsapp-bot/logger"
	"whatsapp-bot/plugin"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	qrterminal "github.com/mdp/qrterminal/v3"
)

// Bot is the main WhatsApp bot instance.
type Bot struct {
	ctx     context.Context
	cancel  context.CancelFunc
	cfg     *config.Config
	log     *logger.Logger
	client  *whatsmeow.Client
	registry *plugin.Registry
}

// New initializes a new Bot with storage and plugin registry.
func New(ctx context.Context, cfg *config.Config, log *logger.Logger) (*Bot, error) {
	ctx, cancel := context.WithCancel(ctx)

	// Setup whatsmeow logger (suppress noise)
	waLogger := waLog.Stdout("WhatsApp", "WARN", true)

	// Setup SQLite store
	container, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", cfg.DBPath), waLogger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	// Load or create device
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	client := whatsmeow.NewClient(deviceStore, waLogger)

	b := &Bot{
		ctx:      ctx,
		cancel:   cancel,
		cfg:      cfg,
		log:      log,
		client:   client,
		registry: plugin.NewRegistry(cfg.Prefix, log),
	}

	// Register all plugins
	b.registerPlugins()

	// Register event handler
	client.AddEventHandler(b.eventHandler)

	return b, nil
}

// Start connects to WhatsApp using QR or pairing code.
func (b *Bot) Start() error {
	if b.client.Store.ID == nil {
		return b.firstLogin()
	}

	b.log.Info("Resuming existing session", "jid", b.client.Store.ID)
	return b.connect()
}

// Stop disconnects the client cleanly.
func (b *Bot) Stop() {
	b.cancel()
	b.client.Disconnect()
}

// ─── Connection ──────────────────────────────────────────────────────────────

func (b *Bot) firstLogin() error {
	switch b.cfg.ConnectMethod {
	case config.ConnectQR:
		return b.loginWithQR()
	case config.ConnectPairingCode:
		return b.loginWithPairingCode()
	default:
		return fmt.Errorf("unknown connect method: %s", b.cfg.ConnectMethod)
	}
}

func (b *Bot) loginWithQR() error {
	b.log.Info("Connecting via QR code...")

	qrChan, err := b.client.GetQRChannel(b.ctx)
	if err != nil {
		if b.client.Store.ID != nil {
			return b.connect()
		}
		return fmt.Errorf("failed to get QR channel: %w", err)
	}

	if err := b.client.Connect(); err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	for evt := range qrChan {
		switch evt.Event {
		case "code":
			b.log.Info("Scan this QR code with your WhatsApp app:")
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		case "success":
			b.log.Info("✅ QR login successful!")
		case "timeout":
			return fmt.Errorf("QR code timed out — please restart")
		case "error":
			return fmt.Errorf("QR error: %v", evt.Error)
		}
	}
	return nil
}

func (b *Bot) loginWithPairingCode() error {
	b.log.Info("Connecting via pairing code...", "phone", b.cfg.PhoneNumber)

	if err := b.client.Connect(); err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	// Small delay to let the connection stabilize
	time.Sleep(2 * time.Second)

	code, err := b.client.PairPhone(
		b.cfg.PhoneNumber,
		true,
		whatsmeow.PairClientChrome,
		"Chrome (Linux)",
	)
	if err != nil {
		return fmt.Errorf("pairing failed: %w", err)
	}

	b.log.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	b.log.Info(fmt.Sprintf("🔑 Pairing Code: %s", code))
	b.log.Info("Open WhatsApp → Linked Devices → Link a Device → Enter code above")
	b.log.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Wait for pairing to complete
	<-b.ctx.Done()
	return nil
}

func (b *Bot) connect() error {
	if err := b.client.Connect(); err != nil {
		return fmt.Errorf("reconnect failed: %w", err)
	}
	b.log.Info("✅ Connected to WhatsApp")
	return nil
}

// ─── Event Handler ────────────────────────────────────────────────────────────

func (b *Bot) eventHandler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {

	case *events.Connected:
		b.log.Info("✅ WhatsApp connected", "jid", b.client.Store.ID)

	case *events.Disconnected:
		b.log.Warn("⚠️  WhatsApp disconnected")
		if b.cfg.AutoReconnect {
			b.scheduleReconnect()
		}

	case *events.LoggedOut:
		b.log.Warn("🚪 Logged out — clearing session")
		_ = b.client.Store.Delete()
		b.Stop()

	case *events.Message:
		b.handleMessage(evt)

	case *events.Receipt:
		// Optional: handle delivery/read receipts

	case *events.HistorySync:
		// Optional: handle history sync

	default:
		// Ignore unhandled events silently
	}
}

func (b *Bot) scheduleReconnect() {
	go func() {
		delay := 5 * time.Second
		for attempt := 1; ; attempt++ {
			select {
			case <-b.ctx.Done():
				return
			case <-time.After(delay):
			}

			b.log.Info("Attempting reconnect...", "attempt", attempt)
			if err := b.connect(); err != nil {
				b.log.Warn("Reconnect failed", "error", err, "retrying_in", delay)
				if delay < 60*time.Second {
					delay *= 2
				}
				continue
			}
			b.log.Info("✅ Reconnected successfully")
			return
		}
	}()
}

// ─── Message Handling ─────────────────────────────────────────────────────────

func (b *Bot) handleMessage(evt *events.Message) {
	// Skip status messages and messages from self
	if evt.Info.IsFromMe || evt.Info.Chat.Server == types.BroadcastServer {
		return
	}

	text := extractText(evt.Message)
	if text == "" {
		return
	}

	b.log.Debug("Incoming message",
		"from", evt.Info.Sender.String(),
		"chat", evt.Info.Chat.String(),
		"text", text,
	)

	// Build context and dispatch to plugin registry
	msgCtx := &plugin.MessageContext{
		Client:  b.client,
		Event:   evt,
		Text:    text,
		Sender:  evt.Info.Sender,
		Chat:    evt.Info.Chat,
		IsGroup: evt.Info.IsGroup,
		Log:     b.log,
		Reply: func(msg string) error {
			return b.sendText(evt.Info.Chat, msg, evt.Info.ID)
		},
		Send: func(msg string) error {
			return b.sendText(evt.Info.Chat, msg, "")
		},
	}

	if err := b.registry.Dispatch(b.ctx, msgCtx); err != nil {
		b.log.Error("Plugin dispatch error", "error", err)
	}
}

func (b *Bot) sendText(to types.JID, text, replyTo string) error {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
		},
	}

	// If replying, use simple text to keep it clean
	if replyTo == "" {
		msg = &waE2E.Message{
			Conversation: proto.String(text),
		}
	}

	_, err := b.client.SendMessage(b.ctx, to, msg)
	if err != nil {
		return fmt.Errorf("send message failed: %w", err)
	}
	return nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func extractText(msg *waE2E.Message) string {
	if msg == nil {
		return ""
	}
	if msg.Conversation != nil {
		return strings.TrimSpace(*msg.Conversation)
	}
	if msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.Text != nil {
		return strings.TrimSpace(*msg.ExtendedTextMessage.Text)
	}
	if msg.ImageMessage != nil && msg.ImageMessage.Caption != nil {
		return strings.TrimSpace(*msg.ImageMessage.Caption)
	}
	if msg.VideoMessage != nil && msg.VideoMessage.Caption != nil {
		return strings.TrimSpace(*msg.VideoMessage.Caption)
	}
	return ""
}
