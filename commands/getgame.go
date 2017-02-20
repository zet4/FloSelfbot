package commands

import "fmt"

// GetGame struct handles GetGame Command
type GetGame struct{}

func (gg *GetGame) message(ctx *Context) {
	em := createEmbed(ctx)
	if CurrentGame != "" {
		em.Description = fmt.Sprintf("Current game is **%s**", CurrentGame)
	} else {
		em.Description = "You are currently playing nothing!"
	}
	ctx.SendEm(em)
}

func (gg *GetGame) description() string { return "Gets your current game." }
func (gg *GetGame) usage() string       { return "" }
func (gg *GetGame) detailed() string {
	return "Because of discord you cant see the change yourself, so why not make a command to see it!"
}
func (gg *GetGame) subcommands() map[string]Command { return make(map[string]Command) }
