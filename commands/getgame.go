package commands

import "fmt"

type GetGame struct{}

func (gg *GetGame) Message(ctx *Context) {
	em := createEmbed(ctx)
	if currentgame != "" {
		em.Description = fmt.Sprintf("Current game is **%s**", currentgame)
	} else {
		em.Description = "You are currently playing nothing!"
	}
	ctx.SendEm(em)
}

func (gg *GetGame) Description() string { return "Gets your current game." }
func (gg *GetGame) Usage() string       { return "" }
func (gg *GetGame) Detailed() string {
	return "Because of discord you cant see the change yourself, so why not make a command to see it!"
}
func (gg *GetGame) Subcommands() map[string]Command { return make(map[string]Command) }
