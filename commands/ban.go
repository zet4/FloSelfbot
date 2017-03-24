package commands

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var modUsageString = "<usermention> [reason]"

// Ban struct handles Ban Command
type Ban struct{}

func (e *Ban) message(ctx *Context) {
	p, _ := ctx.Sess.UserChannelPermissions(ctx.Sess.State.User.ID, ctx.Mess.ChannelID)
	if p&discordgo.PermissionBanMembers != discordgo.PermissionBanMembers {
		ctx.QuickSendEm("You do not have permission to ban members!")
		return
	}
	if len(ctx.Args) != 0 {
		var reason string
		if len(ctx.Mess.Mentions) < 1 {
			ctx.QuickSendEm("You didnt mention a user!")
			return
		}
		reason = strings.TrimSpace(regexp.MustCompile(`^(.*?)<@!?\d+>(.*)$`).ReplaceAllString(ctx.Argstr, "${1}$2"))

		err := ctx.Sess.GuildBanCreate(ctx.Guild.ID, ctx.Mess.Mentions[0].ID, 0)
		if err != nil {
			ctx.QuickSendEm("Error banning user: " + err.Error())
			return
		}

		em := createEmbed(ctx)
		em.Author = &discordgo.MessageEmbedAuthor{IconURL: discordgo.EndpointUserAvatar(ctx.Mess.Mentions[0].ID, ctx.Mess.Mentions[0].Avatar), Name: fmt.Sprintf("Banned: %s#%s (%s)", ctx.Mess.Mentions[0].Username, ctx.Mess.Mentions[0].Discriminator, ctx.Mess.Mentions[0].ID)}
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

func (e *Ban) description() string { return "Bans user from your server" }
func (e *Ban) usage() string       { return modUsageString }
func (e *Ban) detailed() string {
	return "Bans a user from your server, you can add a reason for your ban if desired."
}
func (e *Ban) subcommands() map[string]Command { return make(map[string]Command) }
