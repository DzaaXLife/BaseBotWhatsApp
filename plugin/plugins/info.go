package plugins

import (
	"context"
	"fmt"
	"time"

	"whatsapp-bot/plugin"
)

// InfoPlugin shows info about the current chat.
type InfoPlugin struct{}

func NewInfoPlugin() *InfoPlugin { return &InfoPlugin{} }

func (p *InfoPlugin) Name() string        { return "info" }
func (p *InfoPlugin) Description() string { return "Tampilkan info chat dan pengirim" }
func (p *InfoPlugin) Commands() []string  { return []string{"info", "whoami"} }

func (p *InfoPlugin) Execute(ctx context.Context, msg *plugin.MessageContext) error {
	chatType := "Private"
	if msg.IsGroup {
		chatType = "Group"
	}

	text := fmt.Sprintf(
		"ℹ️ *Info*\n"+
			"━━━━━━━━━━━━━━━\n"+
			"👤 Sender: `%s`\n"+
			"💬 Chat: `%s`\n"+
			"🏷️ Type: %s\n"+
			"🕐 Time: %s",
		msg.Sender.String(),
		msg.Chat.String(),
		chatType,
		time.Now().Format("02 Jan 2006, 15:04:05"),
	)

	return msg.Reply(text)
}
