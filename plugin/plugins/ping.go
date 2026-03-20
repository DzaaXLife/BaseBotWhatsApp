package plugins

import (
	"context"
	"fmt"
	"time"

	"whatsapp-bot/plugin"
)

// PingPlugin responds with latency info.
type PingPlugin struct{}

func NewPingPlugin() *PingPlugin { return &PingPlugin{} }

func (p *PingPlugin) Name() string        { return "ping" }
func (p *PingPlugin) Description() string { return "Cek apakah bot aktif dan responsif" }
func (p *PingPlugin) Commands() []string  { return []string{"ping"} }

func (p *PingPlugin) Execute(ctx context.Context, msg *plugin.MessageContext) error {
	start := time.Now()
	latency := time.Since(start)
	return msg.Reply(fmt.Sprintf("🏓 *Pong!*\nLatency: %dms", latency.Milliseconds()))
}
