package bot

import (
	"whatsapp-bot/plugin"
	"whatsapp-bot/plugin/plugins"
)

// registerPlugins wires all plugins into the registry.
// Add new plugins here.
func (b *Bot) registerPlugins() {
	r := b.registry

	// ── Built-in plugins ──────────────────────────────────────────
	r.Register(plugins.NewPingPlugin())
	r.Register(plugins.NewHelpPlugin(r))
	r.Register(plugins.NewEchoPlugin())
	r.Register(plugins.NewInfoPlugin())

	// ── Add your custom plugins below ─────────────────────────────
	// r.Register(plugins.NewYourPlugin())
}
