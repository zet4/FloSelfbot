package commands

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

// React struct handles React Command
type React struct{}

func (r *React) message(ctx *Context) {
	if len(ctx.Args) < 2 {
		ctx.QuickSendEm("Not enough arguments passed!")
		return
	}
	var qmess *discordgo.Message
	msgs, err := ctx.Sess.ChannelMessages(ctx.Mess.ChannelID, 3, ctx.Mess.ID, "", ctx.Args[0])
	logerror(err)
	for _, msg := range msgs {
		if msg.ID == ctx.Args[0] {
			qmess = msg
		}
	}
	if qmess == nil {
		ctx.QuickSendEm("Message not found!")
		return
	}
	if len(ctx.Args[1:]) > 20-len(qmess.Reactions) {
		ctx.QuickSendEm("Too many emojis to add!")
		return
	}
	for _, e := range ctx.Args[1:] {
		e = regexp.MustCompile(`<:(.*?):(\d+)>`).ReplaceAllString(e, "$1:$2")
		logerror(ctx.Sess.MessageReactionAdd(qmess.ChannelID, qmess.ID, e))
	}
}

func (r *React) description() string             { return "Reacts to messageID with all emojis after the ID" }
func (r *React) usage() string                   { return "<MessageID> [emoji1] [emoji2]..." }
func (r *React) detailed() string                { return "Reacts to messageID with all emojis after the ID" }
func (r *React) subcommands() map[string]Command { return make(map[string]Command) }
