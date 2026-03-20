package plugins

// ─────────────────────────────────────────────────────────────────────────────
// Template Plugin — salin file ini untuk membuat plugin baru
// ─────────────────────────────────────────────────────────────────────────────
//
// Langkah-langkah:
//  1. Salin file ini ke plugin/plugins/nama_plugin.go
//  2. Ganti semua "Template" → nama plugin kamu
//  3. Isi Name(), Description(), Commands(), Execute()
//  4. Daftarkan di bot/plugins.go dengan: r.Register(plugins.NewTemplatePlugin())
//
// ─────────────────────────────────────────────────────────────────────────────

import (
	"context"
	"fmt"

	"whatsapp-bot/plugin"
)

// TemplatePlugin is a starter template for new plugins.
type TemplatePlugin struct {
	// Add any dependencies here (e.g., http.Client, DB handle, etc.)
}

// NewTemplatePlugin creates a new TemplatePlugin.
func NewTemplatePlugin() *TemplatePlugin {
	return &TemplatePlugin{}
}

// Name returns a unique identifier for this plugin.
func (p *TemplatePlugin) Name() string { return "template" }

// Description is shown in the !help command.
func (p *TemplatePlugin) Description() string {
	return "Contoh plugin — ganti dengan deskripsi plugin kamu"
}

// Commands lists all command words this plugin handles (without prefix).
// The first entry is the primary command; the rest are aliases.
func (p *TemplatePlugin) Commands() []string {
	return []string{"template", "tmpl"}
}

// Execute runs when a matching command is received.
//
// msg.Args[0]  = command word (e.g. "template")
// msg.Args[1:] = arguments passed by the user
//
// Use msg.Reply() to quote the user's message.
// Use msg.Send()  to send without quoting.
// Return a non-nil error to trigger automatic error notification.
func (p *TemplatePlugin) Execute(ctx context.Context, msg *plugin.MessageContext) error {
	// ── Argument validation ───────────────────────────────────────
	if len(msg.Args) < 2 {
		return fmt.Errorf("penggunaan: !template <argumen>")
	}

	// ── Your logic here ──────────────────────────────────────────
	arg := msg.Args[1]
	result := fmt.Sprintf("📦 Kamu mengirim: *%s*", arg)

	// ── Owner-only guard example ──────────────────────────────────
	// ownerJID := "6281234567890@s.whatsapp.net"
	// if msg.Sender.String() != ownerJID {
	// 	return fmt.Errorf("perintah ini hanya untuk owner")
	// }

	// ── Group-only guard example ──────────────────────────────────
	// if !msg.IsGroup {
	// 	return fmt.Errorf("perintah ini hanya bisa digunakan di grup")
	// }

	return msg.Reply(result)
}
