package commands

import "github.com/bwmarrin/discordgo"

// Status struct handles Status Command
type Status struct{}

func (s *Status) message(ctx *Context) {
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

func (s *Status) description() string             { return "Sets your status. (online|invisible|dnd|idle)" }
func (s *Status) usage() string                   { return "<online|invisible|dnd|idle>" }
func (s *Status) detailed() string                { return "Sets your status. (online|invisible|dnd|idle)" }
func (s *Status) subcommands() map[string]Command { return make(map[string]Command) }
