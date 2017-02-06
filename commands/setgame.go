package commands

import (
	"fmt"
	"strings"
)

// SetGame struct handles SetGame Command
type SetGame struct{}

func (sg *SetGame) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		game := strings.Join(ctx.Args, " ")
		em.Description = fmt.Sprintf("Changed game to **%s**", game)
		ctx.Sess.UpdateStatus(0, game)
		currentgame = game
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "No game name specified!"
		ctx.SendEm(em)
	}
}

func (sg *SetGame) description() string { return "Sets your game to anything you like" }
func (sg *SetGame) usage() string       { return "<game>" }
func (sg *SetGame) detailed() string {
	return "Changes your 'Playing' status on discord (Because of discord you cant see the change yourself.)"
}
func (sg *SetGame) subcommands() map[string]Command { return make(map[string]Command) }
