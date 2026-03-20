package plugins

import (
	"context"
	"fmt"
	"strings"

	"whatsapp-bot/plugin"
)

// HelpPlugin dynamically lists all registered commands.
type HelpPlugin struct {
	registry *plugin.Registry
}

func NewHelpPlugin(r *plugin.Registry) *HelpPlugin {
	return &HelpPlugin{registry: r}
}

func (p *HelpPlugin) Name() string        { return "help" }
func (p *HelpPlugin) Description() string { return "Tampilkan daftar semua perintah" }
func (p *HelpPlugin) Commands() []string  { return []string{"help", "menu", "?"} }

func (p *HelpPlugin) Execute(ctx context.Context, msg *plugin.MessageContext) error {
	var sb strings.Builder
	sb.WriteString("🤖 *Daftar Perintah Bot*\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━\n\n")

	for _, pl := range p.registry.List() {
		cmds := pl.Commands()
		primary := cmds[0]
		aliases := ""
		if len(cmds) > 1 {
			aliases = fmt.Sprintf(" _(alias: %s)_", strings.Join(cmds[1:], ", "))
		}
		sb.WriteString(fmt.Sprintf("• *!%s*%s\n  %s\n\n", primary, aliases, pl.Description()))
	}

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString("Gunakan prefix `!` sebelum perintah.")

	return msg.Reply(sb.String())
}
