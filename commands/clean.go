package commands

import "strconv"

// Clean struct handles Clean Command
type Clean struct{}

func (c *Clean) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		var mID string
		limit, err := strconv.Atoi(ctx.Argstr)
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
			if msg.Author.ID == ctx.Sess.State.User.ID {
				err = ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, msg.ID)
				logerror(err)
			}
		}
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify an amount"
		ctx.SendEm(em)
	}
}

func (c *Clean) description() string { return "Cleans up your messages" }
func (c *Clean) usage() string       { return "<amount>/<messageID>" }
func (c *Clean) detailed() string {
	return "If you realise you have been spamming a little, this is the command to use then.\nIf message ID is specified, delete everything you posted until that message. (This does not include that message)"
}
func (c *Clean) subcommands() map[string]Command { return make(map[string]Command) }
