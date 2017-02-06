package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)
// Quote struct handles Quote Command
type Quote struct{}

func (q *Quote) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		var qmess *discordgo.Message
		var mID, cID string
		if len(ctx.Args) > 1 {
			cID = ctx.Args[0]
			mID = ctx.Args[1]
		} else {
			mID = ctx.Args[0]
			cID = ctx.Mess.ChannelID
		}
		msgs, err := ctx.Sess.ChannelMessages(cID, 3, ctx.Mess.ID, "", mID)
		for _, msg := range msgs {
			if msg.ID == mID {
				qmess = msg
			}
		}
		if qmess == nil {
			em := createEmbed(ctx)
			em.Description = "message not found"
			ctx.SendEm(em)
			return
		}
		emauthor := &discordgo.MessageEmbedAuthor{Name: qmess.Author.Username, IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", qmess.Author.ID, qmess.Author.Avatar)}
		timestamp, err := qmess.Timestamp.Parse()
		logerror(err)
		timestampo := timestamp.Local().Format(time.ANSIC)
		emfooter := &discordgo.MessageEmbedFooter{Text: "Sent | " + timestampo}
		emcolor := ctx.Sess.State.UserColor(qmess.Author.ID, qmess.ChannelID)
		em := &discordgo.MessageEmbed{Author: emauthor, Footer: emfooter, Description: qmess.Content, Color: emcolor}
		ctx.SendEmNoDelete(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify a message ID"
		ctx.SendEm(em)
	}
}

func (q *Quote) description() string { return "Quotes a message" }
func (q *Quote) usage() string       { return "<messageID> or <channelID> <messageID>" }
func (q *Quote) detailed() string {
	return "To find messageID and channelID you first need to turn on Developer mode in discord, then right click any message/channel and click 'Copy ID'"
}
func (q *Quote) subcommands() map[string]Command { return make(map[string]Command) }
