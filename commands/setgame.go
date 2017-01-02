package commands

import (
	"fmt"
	"strings"
)

type SetGame struct{}

func (sg *SetGame) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		game := strings.Join(ctx.Args, " ")
		em.Description = fmt.Sprintf("Changed game to **%s**", game)
		ctx.Sess.UpdateStatus(0, game)
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "No game name specified!"
		ctx.SendEm(em)
	}
}

func (sg *SetGame) Description() string { return "Sets your game to anything you like" }
func (sg *SetGame) Usage() string       { return "<game>" }
func (sg *SetGame) Detailed() string {
	return "Changes your 'Playing' status on discord (Because of discord you cant see the change yourself.)"
}
func (sg *SetGame) Subcommands() map[string]Command { return make(map[string]Command) }
