package commands

import (
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
)

// MessageCache stores messages using the ID as key
var MessageCache = cache.New(10*time.Minute, 20*time.Minute)

// QuoteRegex struct handles QuoteRegex Subcommand
type QuoteRegex struct{}

func (q *QuoteRegex) message(ctx *Context) {
	var qmess *discordgo.Message

	regex, err := regexp.Compile(ctx.Argstr)

	if err != nil {
		ctx.QuickSendEm("Error compiling your regex!")
	}

	msgs, _ := ctx.Sess.ChannelMessages(ctx.Mess.ChannelID, 100, ctx.Mess.ID, "", "")
	for _, msg := range msgs {
		if len(regex.FindStringSubmatch(msg.Content)) > 0 {
			qmess = msg
			break
		}
	}
	if qmess == nil {
		em := createEmbed(ctx)
		em.Description = "Match not found!"
		ctx.SendEm(em)
		return
	}

	// var guild *discordgo.Guild
	var authorIcon, guildIcon string

	if !ctx.Channel.IsPrivate {
		if len(ctx.Guild.Icon) > 0 {
			guildIcon = discordgo.EndpointGuildIcon(ctx.Guild.ID, ctx.Guild.Icon)
		}
	}

	authorIcon = discordgo.EndpointUserAvatar(qmess.Author.ID, qmess.Author.Avatar)

	emauthor := &discordgo.MessageEmbedAuthor{Name: qmess.Author.Username, IconURL: authorIcon}
	timestamp, err := qmess.Timestamp.Parse()
	logerror(err)
	timestampo := timestamp.Local().Format(time.ANSIC)
	emfooter := &discordgo.MessageEmbedFooter{Text: "Sent | " + timestampo, IconURL: guildIcon}
	emcolor := ctx.Sess.State.UserColor(qmess.Author.ID, qmess.ChannelID)
	em := &discordgo.MessageEmbed{Author: emauthor, Footer: emfooter, Description: qmess.Content, Color: emcolor}
	ctx.SendEmNoDelete(em)
}

func (q *QuoteRegex) description() string {
	return "Quotes a message in your channel with regex (only the last 100 msgs)"
}
func (q *QuoteRegex) usage() string { return "<regex>" }
func (q *QuoteRegex) detailed() string {
	return "Use regex to find the latest match"
}
func (q *QuoteRegex) subcommands() map[string]Command { return make(map[string]Command) }

// Quote struct handles Quote Command
type Quote struct{}

func (q *Quote) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		var qmess *discordgo.Message
		var mID, cID string
		var ch *discordgo.Channel
		if len(ctx.Args) > 1 {
			cID = ctx.Args[0]
			mID = ctx.Args[1]
		} else {
			mID = ctx.Args[0]
			cID = ctx.Mess.ChannelID
		}
		ch, err := ctx.Sess.State.Channel(cID)
		if err != nil {
			chs, _ := ctx.Sess.UserChannels()
			for _, c := range chs {
				if c.Recipient.ID == cID {
					ch = c
					break
				}
			}
		}
		msgs, _ := ctx.Sess.ChannelMessages(ch.ID, 3, ctx.Mess.ID, "", mID)
		for _, msg := range msgs {
			if msg.ID == mID {
				qmess = msg
			}
		}
		if qmess == nil {
			if x, found := MessageCache.Get(mID); found {
				qmess = x.(*discordgo.Message)
				cID = qmess.ChannelID
			} else {
				em := createEmbed(ctx)
				em.Description = "Message not found!"
				ctx.SendEm(em)
				return
			}
		}

		// var guild *discordgo.Guild
		var authorIcon, guildIcon string

		if err == nil && !ch.IsPrivate {
			guild, _ := ctx.Sess.State.Guild(ch.GuildID)
			if len(guild.Icon) > 0 {
				guildIcon = discordgo.EndpointGuildIcon(guild.ID, guild.Icon)
			}
		}

		authorIcon = discordgo.EndpointUserAvatar(qmess.Author.ID, qmess.Author.Avatar)

		emauthor := &discordgo.MessageEmbedAuthor{Name: qmess.Author.Username, IconURL: authorIcon}
		timestamp, err := qmess.Timestamp.Parse()
		logerror(err)
		timestampo := timestamp.Local().Format(time.ANSIC)
		emfooter := &discordgo.MessageEmbedFooter{Text: "Sent | " + timestampo, IconURL: guildIcon}
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
func (q *Quote) subcommands() map[string]Command { return map[string]Command{"regex": &QuoteRegex{}} }
