package plugin

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"whatsapp-bot/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// ─── Interfaces & Types ───────────────────────────────────────────────────────

// Plugin is the interface every plugin must implement.
type Plugin interface {
	// Name returns the unique plugin identifier.
	Name() string
	// Description is shown in the help command.
	Description() string
	// Commands returns all commands this plugin handles (without prefix).
	Commands() []string
	// Execute is called when a matching command is received.
	Execute(ctx context.Context, msg *MessageContext) error
}

// MessageContext carries everything a plugin needs to respond.
type MessageContext struct {
	Client  *whatsmeow.Client
	Event   *events.Message
	Text    string // Full raw text
	Args    []string // Command args (split by space, first element is command)
	Sender  types.JID
	Chat    types.JID
	IsGroup bool
	Log     *logger.Logger

	// Helpers
	Reply func(msg string) error // Reply mentioning the message
	Send  func(msg string) error // Send without quoting
}

// Command returns the first arg (the matched command word, without prefix).
func (m *MessageContext) Command() string {
	if len(m.Args) == 0 {
		return ""
	}
	return m.Args[0]
}

// ─── Registry ─────────────────────────────────────────────────────────────────

// Registry holds all registered plugins and routes messages to them.
type Registry struct {
	prefix  string
	log     *logger.Logger
	plugins map[string]Plugin  // command → plugin
	list    []Plugin           // ordered list (for help)
}

func NewRegistry(prefix string, log *logger.Logger) *Registry {
	return &Registry{
		prefix:  prefix,
		log:     log,
		plugins: make(map[string]Plugin),
	}
}

// Register adds a plugin to the registry.
// Panics on duplicate command registration (caught at startup).
func (r *Registry) Register(p Plugin) {
	for _, cmd := range p.Commands() {
		key := strings.ToLower(cmd)
		if existing, ok := r.plugins[key]; ok {
			panic(fmt.Sprintf(
				"plugin conflict: command %q already registered by %q, cannot register %q",
				cmd, existing.Name(), p.Name(),
			))
		}
		r.plugins[key] = p
	}
	r.list = append(r.list, p)
	r.log.Info("Plugin registered", "name", p.Name(), "commands", p.Commands())
}

// List returns all registered plugins in registration order.
func (r *Registry) List() []Plugin {
	return r.list
}

// Dispatch routes a message to the correct plugin.
// Returns nil if no command matched.
func (r *Registry) Dispatch(ctx context.Context, msg *MessageContext) (err error) {
	// Must start with prefix
	if !strings.HasPrefix(msg.Text, r.prefix) {
		return nil
	}

	// Parse command and args
	trimmed := strings.TrimPrefix(msg.Text, r.prefix)
	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return nil
	}

	cmd := strings.ToLower(parts[0])
	msg.Args = parts // Args[0] = command word (without prefix), Args[1:] = arguments

	p, ok := r.plugins[cmd]
	if !ok {
		return nil // Unknown command — silently ignore
	}

	r.log.Info("Dispatching command",
		"plugin", p.Name(),
		"command", cmd,
		"sender", msg.Sender.String(),
	)

	// Execute with panic recovery
	return safeExecute(ctx, p, msg)
}

// safeExecute runs a plugin and recovers from panics.
func safeExecute(ctx context.Context, p Plugin, msg *MessageContext) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			err = fmt.Errorf("panic in plugin %q: %v\n%s", p.Name(), r, stack)
			// Try to notify user
			_ = msg.Reply("⚠️ Terjadi error internal pada perintah ini.")
		}
	}()

	if execErr := p.Execute(ctx, msg); execErr != nil {
		// Notify user about the error
		_ = msg.Reply(fmt.Sprintf("❌ Error: %s", execErr.Error()))
		return fmt.Errorf("plugin %q execute error: %w", p.Name(), execErr)
	}
	return nil
}
