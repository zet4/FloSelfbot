package commands

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Kick struct handles Kick Command
type Kick struct{}

func (e *Kick) message(ctx *Context) {
	p, _ := ctx.Sess.UserChannelPermissions(ctx.Sess.State.User.ID, ctx.Mess.ChannelID)
	if p&discordgo.PermissionKickMembers != discordgo.PermissionKickMembers {
		ctx.QuickSendEm("You do not have permission to kick members!")
		return
	}
	if len(ctx.Args) != 0 {
		if len(ctx.Mess.Mentions) < 1 {
			ctx.QuickSendEm("You didnt mention a user!")
			return
		}
		var reason string
		reason = strings.TrimSpace(regexp.MustCompile(`^(.*?)<@!?\d+>(.*)$`).ReplaceAllString(ctx.Argstr, "${1}$2"))

		err := ctx.Sess.GuildMemberDelete(ctx.Guild.ID, ctx.Mess.Mentions[0].ID)
		if err != nil {
			ctx.QuickSendEm("Error kicking user: " + err.Error())
			return
		}

		em := createEmbed(ctx)
		em.Author = &discordgo.MessageEmbedAuthor{IconURL: discordgo.EndpointUserAvatar(ctx.Mess.Mentions[0].ID, ctx.Mess.Mentions[0].Avatar), Name: fmt.Sprintf("Kicked: %s#%s (%s)", ctx.Mess.Mentions[0].Username, ctx.Mess.Mentions[0].Discriminator, ctx.Mess.Mentions[0].ID)}
		if reason != "" {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Reason", Value: reason, Inline: true})
		}
		timestamp := time.Now().UTC().Format("2006-01-02 15:04:05") + " UTC"
		em.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("By: %s#%s (%s) | %s", ctx.Mess.Author.Username, ctx.Mess.Author.Discriminator, ctx.Mess.Author.ID, timestamp), IconURL: discordgo.EndpointUserAvatar(ctx.Mess.Author.ID, ctx.Mess.Author.Avatar)}
		ctx.SendEmNoDelete(em)
	} else {
		ctx.QuickSendEm("No user specified!")
	}
}

func (e *Kick) description() string { return "Kicks user from your server" }
func (e *Kick) usage() string       { return modUsageString }
func (e *Kick) detailed() string {
	return "Kicks a user from your server, you can add a reason for your kick if desired."
}
func (e *Kick) subcommands() map[string]Command { return make(map[string]Command) }
