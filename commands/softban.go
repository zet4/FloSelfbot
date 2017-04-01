package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Softban struct handles Softban Command
type Softban struct{}

func (e *Softban) message(ctx *Context) {
	p, _ := ctx.Sess.UserChannelPermissions(ctx.Sess.State.User.ID, ctx.Mess.ChannelID)
	if p&discordgo.PermissionBanMembers != discordgo.PermissionBanMembers {
		ctx.QuickSendEm("You do not have permission to ban members!")
		return
	}
	if len(ctx.Args) != 0 {
		var reason string
		split := strings.SplitN(ctx.Argstr, " for ", 2)
		if len(split) > 1 {
			reason = split[1]
		}
		user, err := ctx.GetUser(split[0], ctx.Guild.ID)
		if err != nil {
			return
		}
		err = ctx.Sess.GuildMemberDelete(ctx.Guild.ID, user.ID)
		if err != nil {
			ctx.QuickSendEm("Error kicking user: " + err.Error())
			return
		}

		em := createEmbed(ctx)
		em.Author = &discordgo.MessageEmbedAuthor{IconURL: discordgo.EndpointUserAvatar(user.ID, user.Avatar), Name: fmt.Sprintf("Softbanned: %s#%s (%s)", user.Username, user.Discriminator, user.ID)}
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

func (e *Softban) description() string { return "Softbans user from your server" }
func (e *Softban) usage() string       { return modUsageString }
func (e *Softban) detailed() string {
	return "Softbans a user on your server (ban, removing 1 day, unban), you can add a reason for your softban if desired."
}
func (e *Softban) subcommands() map[string]Command { return make(map[string]Command) }
