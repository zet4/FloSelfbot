package commands

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Sinfo struct handles Sinfo Command
type Sinfo struct{}

func (s *Sinfo) message(ctx *Context) {
	em := createEmbed(ctx)
	em.Author = &discordgo.MessageEmbedAuthor{
		Name:    fmt.Sprintf("Server Info: %s", ctx.Guild.Name),
		IconURL: "https://twemoji.maxcdn.com/36x36/2139.png",
	}
	if ctx.Guild.Icon != "" {
		em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.jpg", ctx.Guild.ID, ctx.Guild.Icon)}
	}
	em.Fields = make([]*discordgo.MessageEmbedField, 0)
	if len(ctx.Guild.Emojis) > 0 {
		var buf bytes.Buffer
		for _, emote := range ctx.Guild.Emojis[0:int(math.Min(float64(40), float64(len(ctx.Guild.Emojis))))] {
			fmt.Fprintf(&buf, "<:x:%s>", emote.ID)
		}
		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:  "Emojis (up to 40)",
			Value: buf.String(),
		})
	}
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
		Name:   "Owner",
		Value:  "<@" + ctx.Guild.OwnerID + ">",
		Inline: true,
	})
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
		Name:   "Roles",
		Value:  strconv.Itoa(len(ctx.Guild.Roles)),
		Inline: true,
	})
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
		Name:   "Emojis",
		Value:  strconv.Itoa(len(ctx.Guild.Emojis)),
		Inline: true,
	})
	var regular, voice, hidden int
	for _, ch := range ctx.Guild.Channels {
		regular++
		if ch.Type == "voice" {
			voice++
			continue
		}
		perms, err := ctx.Sess.State.UserChannelPermissions(ctx.Sess.State.User.ID, ch.ID)
		if err == nil && perms&discordgo.PermissionReadMessages == 0 {
			hidden++
		}
	}
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
		Name:   "Channels",
		Value:  fmt.Sprintf("%d channels (%d voice, %d hidden)", regular, voice, hidden),
		Inline: true,
	})
	var bots int
	for _, mb := range ctx.Guild.Members {
		if mb.User.Bot {
			bots++
		}
	}
	var botmsg string
	if bots > 0 {
		botmsg = fmt.Sprintf(", %d or %.2f%% of members are bots", bots, (float64(bots)*float64(100))/float64(len(ctx.Guild.Members)))
	}
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
		Name:   "Members",
		Value:  fmt.Sprintf("%d total%s", len(ctx.Guild.Members), botmsg),
		Inline: true,
	})
	ctx.SendEm(em)
}

func (s *Sinfo) description() string             { return "Server info" }
func (s *Sinfo) usage() string                   { return "" }
func (s *Sinfo) detailed() string                { return "Returns server info for the server you are currently in." }
func (s *Sinfo) subcommands() map[string]Command { return make(map[string]Command) }
