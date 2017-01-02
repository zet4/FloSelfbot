package commands

import "strconv"

type Clean struct{}

func (c *Clean) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		limit, err := strconv.Atoi(ctx.Args[0])
		logerror(err)
		msgs, err := ctx.Sess.ChannelMessages(ctx.Mess.ChannelID, limit, ctx.Mess.ID, "")
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

func (c *Clean) Description() string { return "Cleans up your messages" }
func (c *Clean) Usage() string       { return "<amount>" }
func (c *Clean) Detailed() string {
	return "If you realise you have been spamming a little, this is the command to use then."
}
func (c *Clean) Subcommands() map[string]Command { return make(map[string]Command) }
