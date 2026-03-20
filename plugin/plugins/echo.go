package plugins

import (
	"context"
	"fmt"
	"strings"

	"whatsapp-bot/plugin"
)

// EchoPlugin echoes back what the user says.
type EchoPlugin struct{}

func NewEchoPlugin() *EchoPlugin { return &EchoPlugin{} }

func (p *EchoPlugin) Name() string        { return "echo" }
func (p *EchoPlugin) Description() string { return "Ulangi pesan yang kamu kirim (!echo <pesan>)" }
func (p *EchoPlugin) Commands() []string  { return []string{"echo", "say"} }

func (p *EchoPlugin) Execute(ctx context.Context, msg *plugin.MessageContext) error {
	if len(msg.Args) < 2 {
		return fmt.Errorf("penggunaan: !echo <pesan>")
	}

	text := strings.Join(msg.Args[1:], " ")
	return msg.Send(text)
}
