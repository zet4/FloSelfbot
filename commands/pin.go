package commands

type Pin struct{}

func (p *Pin) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		err := ctx.Sess.ChannelMessagePin(ctx.Mess.ChannelID, ctx.Args[0])
		logerror(err)
		if err != nil {
			em := createEmbed(ctx)
			em.Description = "Error: " + err.Error()
			ctx.SendEm(em)
		} else {
			em := createEmbed(ctx)
			em.Description = "Sucessfully pinned message `" + ctx.Args[0] + "`"
			ctx.SendEm(em)
		}

	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify a message ID"
		ctx.SendEm(em)
	}
}

func (p *Pin) Description() string { return "Pins a message" }
func (p *Pin) Usage() string       { return "<messageID>" }
func (p *Pin) Detailed() string {
	return "To find messageIDyou first need to turn on Developer mode in discord, then right click any message and click 'Copy ID'"
}
func (p *Pin) Subcommands() map[string]Command { return make(map[string]Command) }
