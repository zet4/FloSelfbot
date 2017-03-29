package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Uinfo struct handles Uinfo Command
type Uinfo struct{}

func (u *Uinfo) message(ctx *Context) {
	var user *discordgo.User
	var users []*discordgo.User
	var err error

	if len(ctx.Args) > 0 {
		if !ctx.Channel.IsPrivate {
			users, err = ctx.GuildGetUserByName(ctx.Argstr, ctx.Channel.GuildID)
			if err != nil {
				ctx.QuickSendEm("Error collecting users: " + err.Error())
				return
			}
			if len(users) < 1 {
				users, err = ctx.GetUserByName(ctx.Argstr)
				if len(users) < 1 {
					ctx.QuickSendEm("No user found with name **" + ctx.Argstr + "**")
				}
				if err != nil {
					ctx.QuickSendEm("Error collecting users: " + err.Error())
					return
				}
			}
			if len(users) > 1 {
				ctx.ParseTooManyUsers(ctx.Argstr, users)
				return
			}
			user = users[0]
		} else {
			users, err = ctx.GetUserByName(ctx.Argstr)
			if err != nil {
				ctx.QuickSendEm("Error collecting users: " + err.Error())
				return
			}
			if len(users) < 1 {
				ctx.QuickSendEm("No user found with name **" + ctx.Argstr + "**")
			}
			if len(users) > 1 {
				ctx.ParseTooManyUsers(ctx.Argstr, users)
				return
			}
			user = users[0]
		}
		if len(users) > 1 {
			ctx.ParseTooManyUsers(ctx.Argstr, users)
			return
		}
	} else {
		user = ctx.Mess.Author
	}

	em := createEmbed(ctx)
	em.Author = &discordgo.MessageEmbedAuthor{
		Name:    fmt.Sprintf("User Info: %s#%s", user.Username, user.Discriminator),
		IconURL: "https://twemoji.maxcdn.com/36x36/2139.png",
	}
	em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: discordgo.EndpointUserAvatar(user.ID, user.Avatar)}
	em.Fields = make([]*discordgo.MessageEmbedField, 0)
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
		Name:   "ID",
		Value:  user.ID,
		Inline: true,
	})
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
		Name:   "Mention",
		Value:  "<@" + user.ID + ">",
		Inline: true,
	})
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
		Name:   "Bot",
		Value:  strconv.FormatBool(user.Bot),
		Inline: true,
	})
	if !ctx.Channel.IsPrivate {
		var member *discordgo.Member
		for _, m := range ctx.Guild.Members {
			if m.User.ID == user.ID {
				member = m
				break
			}
		}
		if member != nil {
			if len(member.Roles) > 0 {
				var buf bytes.Buffer
				for _, rid := range member.Roles {
					role, err := ctx.Sess.State.Role(ctx.Guild.ID, rid)
					if err != nil {
						ctx.QuickSendEm("Error getting role: " + err.Error())
					}
					fmt.Fprintf(&buf, "%s, ", role.Name)
				}
				em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
					Name:  "Roles",
					Value: buf.String()[:len(buf.String())-2],
				})
			}
			if t, err := discordgo.Timestamp(member.JoinedAt).Parse(); err == nil {
				em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
					Name:  "Join date",
					Value: fmt.Sprintf("%s (%.2f days ago)", t.Format(time.ANSIC), time.Now().Sub(t).Hours()/24),
				})
			}
		}
	}
	if t, err := ctx.GetCreationTime(user.ID); err == nil {
		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   "Creation",
			Value:  fmt.Sprintf("%s (%.2f days ago)", t.Format(time.ANSIC), time.Now().Sub(t).Hours()/24),
			Inline: true,
		})
	}
	ctx.SendEm(em)
}

func (u *Uinfo) description() string             { return "User info" }
func (u *Uinfo) usage() string                   { return "<user>" }
func (u *Uinfo) detailed() string                { return "Returns user info for the user specified." }
func (u *Uinfo) subcommands() map[string]Command { return make(map[string]Command) }
