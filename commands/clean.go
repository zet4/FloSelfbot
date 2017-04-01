package commands

import (
	"flag"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type bitflag byte

var urlregex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
var inviteregex = regexp.MustCompile(`https:\/\/discord\.gg\/[a-zA-Z0-9]+`)

const (
	all bitflag = 1 << iota
	bots
	links
	invites
	embeds
	attachments
	pins
	invert
)

// Clean struct handles Clean Command
type Clean struct{}

func (c *Clean) message(ctx *Context) {
	if len(ctx.Args) != 0 {

		f := flag.NewFlagSet("CleanSet", flag.ContinueOnError)

		f.Bool("b", false, "Remove messages from bots.")
		f.Bool("l", false, "Remove messages containing links.")
		f.Bool("i", false, "Remove messages containing invites.")
		f.Bool("e", false, "Remove messages containing embeds.")
		f.Bool("a", false, "Remove messages containing attachments.")
		f.Bool("p", false, "Remove messages that are pinned.")
		f.Bool("invert", false, "Reverses the effects of all the flag filters.")

		f.Parse(ctx.Args[1:])

		var flags bitflag

		f.Visit(func(arg *flag.Flag) {
			if arg.Name == "b" {
				flags = flags | bots
			}
			if arg.Name == "l" {
				flags = flags | links
			}
			if arg.Name == "i" {
				flags = flags | invites
			}
			if arg.Name == "e" {
				flags = flags | embeds
			}
			if arg.Name == "a" {
				flags = flags | attachments
			}
			if arg.Name == "p" {
				flags = flags | pins
			}
			if arg.Name == "invert" {
				flags = flags | invert
			}
			flags = flags | all
		})

		var mID string
		limit, err := strconv.Atoi(ctx.Args[0])
		if err != nil {
			ctx.QuickSendEm("Invalid amount specified")
			return
		}
		if limit > 100 {
			mID = ctx.Argstr
			limit = 100
		}
		msgs, err := ctx.Sess.ChannelMessages(ctx.Mess.ChannelID, limit, ctx.Mess.ID, mID, "")
		logerror(err)
		for _, msg := range msgs {
			var delete bool

			if flags&all == 0 && msg.Author.ID == ctx.Sess.State.User.ID {
				delete = true
			}

			if flags&all != 0 {

				p, _ := ctx.Sess.UserChannelPermissions(ctx.Sess.State.User.ID, ctx.Mess.ChannelID)
				if p&discordgo.PermissionManageMessages != discordgo.PermissionManageMessages {
					ctx.QuickSendEm("You do not have permission to clear other user's messages!")
					return
				}

				if flags&invert != 0 {
					flags = ^flags
				}
				if flags&bots != 0 {
					if msg.Author.Bot {
						delete = true
					}
				}
				if flags&links != 0 {
					if len(urlregex.FindStringSubmatch(msg.Content)) > 0 {
						delete = true
					}
				}
				if flags&invites != 0 {
					if len(inviteregex.FindStringSubmatch(msg.Content)) > 0 {
						delete = true
					}
				}
				if flags&embeds != 0 {
					if len(msg.Embeds) > 0 {
						delete = true
					}
				}
				if flags&attachments != 0 {
					if len(msg.Attachments) > 0 {
						delete = true
					}
				}
				if flags&pins != 0 {
					pinmsgs, _ := ctx.Sess.ChannelMessagesPinned(ctx.Mess.ChannelID)
					for _, pm := range pinmsgs {
						if pm.ID == msg.ID {
							delete = true
							break
						}
					}
				}
			}

			if delete {
				err = ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, msg.ID)
				logerror(err)
			}
		}
		err = ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, mID)
		logerror(err)
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify an amount"
		ctx.SendEm(em)
	}
}

func (c *Clean) description() string { return "Cleans up your messages" }
func (c *Clean) usage() string       { return "<amount>/<messageID>" }
func (c *Clean) detailed() string {
	return "If you realise you have been spamming a little, this is the command to use then.\nIf you specify message ID, deletes everything you posted since that message, including that one."
}
func (c *Clean) subcommands() map[string]Command { return make(map[string]Command) }
