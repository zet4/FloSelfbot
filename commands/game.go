package commands

import (
	"fmt"
)

// SetGame struct handles SetGame Subcommand
type SetGame struct{}

func (sg *SetGame) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		game := ctx.Argstr
		em.Description = fmt.Sprintf("Changed game to **%s**", game)
		ctx.Sess.UpdateStatus(0, game)
		CurrentGame = game
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

// GetGame struct handles GetGame Subcommand
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

// Game struct handles Game Command
type Game struct{}

func (c *Game) message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) == 0 {
		em.Description = "Command `game` requires a subcommand!"
	} else {
		em.Description = fmt.Sprintf(`Subcommand "%s" doesn't exist!`, ctx.Args[0])
	}
	ctx.SendEm(em)
}

func (c *Game) description() string { return `Commands for setting and getting your game.` }
func (c *Game) usage() string       { return "" }
func (c *Game) detailed() string {
	return "Commands for setting and getting your game."
}
func (c *Game) subcommands() map[string]Command {
	return map[string]Command{"get": &GetGame{}, "set": &SetGame{}}
}
