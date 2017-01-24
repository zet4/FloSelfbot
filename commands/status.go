package commands

import "github.com/bwmarrin/discordgo"

type Status struct{}

func (s *Status) Message(ctx *Context) {
	if len(ctx.Args) >= 1 {
		_, err := ctx.Sess.UserUpdateStatus(discordgo.Status(ctx.Args[0]))
		if err != nil {
			logerror(err)
			ctx.QuickSendEm("Invalid code entered!")
			return
		}
		ctx.QuickSendEm("Status set to **" + ctx.Args[0] + "**")
	} else {
		ctx.QuickSendEm("No Code specified!")
	}
}

func (s *Status) Description() string             { return "Sets your status. (online|invisible|dnd|idle)" }
func (s *Status) Usage() string                   { return "<online|invisible|dnd|idle>" }
func (s *Status) Detailed() string                { return "Sets your status. (online|invisible|dnd|idle)" }
func (s *Status) Subcommands() map[string]Command { return make(map[string]Command) }
