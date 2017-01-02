package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Quote struct{}

func (q *Quote) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		var qmess *discordgo.Message

		mID := ctx.Args[0]
		msgs, err := ctx.Sess.ChannelMessages(ctx.Mess.ChannelID, 100, ctx.Mess.ID, "")
		for _, msg := range msgs {
			if msg.ID == mID {
				qmess = msg
			}
		}
		if qmess == nil {
			ctx.Sess.ChannelMessageSend(ctx.Mess.ChannelID, "Message not found in last 100 messages.")
		}
		emauthor := &discordgo.MessageEmbedAuthor{Name: qmess.Author.Username, IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", qmess.Author.ID, qmess.Author.Avatar)}
		timestamp, err := qmess.Timestamp.Parse()
		logerror(err)
		timestampo := timestamp.Format(time.ANSIC)
		emfooter := &discordgo.MessageEmbedFooter{Text: "Sent | " + timestampo}
		emcolor := ctx.Sess.State.UserColor(qmess.Author.ID, qmess.ChannelID)
		em := &discordgo.MessageEmbed{Author: emauthor, Footer: emfooter, Description: qmess.Content, Color: emcolor}
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify a message ID"
		ctx.SendEm(em)
	}
}

func (q *Quote) Description() string { return "Quotes a message from the last 100 messages" }
func (q *Quote) Usage() string       { return "<messageID>" }
func (q *Quote) Detailed() string {
	return "To find messageID you first need to turn on Developer mode in discord, then right click any message and click 'Copy ID'"
}
func (q *Quote) Subcommands() map[string]Command { return make(map[string]Command) }
