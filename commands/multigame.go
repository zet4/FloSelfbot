package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Mgtoggle bool

func MultiGameFunc(s *discordgo.Session, conf *Config) {
	for {
		if len(conf.MultiGameStrings) != 0 && conf.MultigameToggled {
			a := conf.MultiGameStrings
			newstring, a := a[0], a[1:]
			conf.MultiGameStrings = append(a, newstring)
			err := s.UpdateStatus(0, newstring)
			currentgame = newstring
			logerror(err)
		}
		if conf.MultiGameMinutes < 1 {
			conf.MultiGameMinutes = 1
			editConfigfile(conf)
		}
		time.Sleep(time.Minute * time.Duration(conf.MultiGameMinutes))
	}
}

type AddMultiGameString struct{}

func (a *AddMultiGameString) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		game := strings.Join(ctx.Args, " ")
		ctx.Conf.MultiGameStrings = append(ctx.Conf.MultiGameStrings, game)
		editConfigfile(ctx.Conf)
		em.Description = fmt.Sprintf("Added **%s** to Multigame", game)
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "No game name specified!"
		ctx.SendEm(em)
	}
}

func (a *AddMultiGameString) Description() string { return "Adds a string to the Multigame" }
func (a *AddMultiGameString) Usage() string       { return "<game>" }
func (a *AddMultiGameString) Detailed() string {
	return "Adds a string to the Multigame."
}
func (a *AddMultiGameString) Subcommands() map[string]Command { return make(map[string]Command) }

type RemoveMultiGameString struct{}

func (r *RemoveMultiGameString) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		var pos int
		em := createEmbed(ctx)
		game := strings.Join(ctx.Args, " ")
		for i, v := range ctx.Conf.MultiGameStrings {
			if game == v {
				pos = i
				break
			} else {
				pos = -1
			}
		}
		if pos == -1 {
			em.Description = "Game `" + game + "` not found in Multigame! (Check for caps and stuff.)"
			ctx.SendEm(em)
		} else {
			ctx.Conf.MultiGameStrings = append(ctx.Conf.MultiGameStrings[:pos], ctx.Conf.MultiGameStrings[pos+1:]...)
			editConfigfile(ctx.Conf)
			em.Description = fmt.Sprintf("Removed **%s** from Multigame", game)
			ctx.SendEm(em)
		}
	} else {
		em := createEmbed(ctx)
		em.Description = "No game name specified!"
		ctx.SendEm(em)
	}
}

func (r *RemoveMultiGameString) Description() string { return "Removes a string from Multigame" }
func (r *RemoveMultiGameString) Usage() string       { return "<game>" }
func (r *RemoveMultiGameString) Detailed() string {
	return "Removes a string from Multigame."
}
func (a *RemoveMultiGameString) Subcommands() map[string]Command { return make(map[string]Command) }

type MultiGameList struct{}

func (l *MultiGameList) Message(ctx *Context) {
	em := createEmbed(ctx)
	desc := "Current strings in Multigame:"
	for _, v := range ctx.Conf.MultiGameStrings {
		desc += fmt.Sprintf("\n`%s`", v)
	}
	em.Description = desc
	ctx.SendEm(em)
}

func (l *MultiGameList) Description() string { return "Returns your Multigame strings" }
func (l *MultiGameList) Usage() string       { return "" }
func (l *MultiGameList) Detailed() string {
	return "Returns the list of all your strings in Multigame"
}
func (l *MultiGameList) Subcommands() map[string]Command { return make(map[string]Command) }

type MultigameTimer struct{}

func (mgt *MultigameTimer) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		delay, err := strconv.Atoi(ctx.Args[0])
		if err != nil {
			em.Description = "Invalid time!"
			ctx.SendEm(em)
			return
		}
		ctx.Conf.MultiGameMinutes = delay
		editConfigfile(ctx.Conf)
		em.Description = fmt.Sprintf("Multigame timer set to **%s** minutes", ctx.Args[0])
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "No int specified!"
		ctx.SendEm(em)
	}
}

func (mgt *MultigameTimer) Description() string {
	return "Changes the delay between switches (in minutes)"
}
func (mgt *MultigameTimer) Usage() string { return "<minutes>" }
func (mgt *MultigameTimer) Detailed() string {
	return "Changes the delay between switches (in minutes)"
}
func (mgt *MultigameTimer) Subcommands() map[string]Command { return make(map[string]Command) }

type MultiGameToggle struct{}

func (mgt *MultiGameToggle) Message(ctx *Context) {
	newtoggle := !ctx.Conf.MultigameToggled
	ctx.Conf.MultigameToggled = newtoggle

	if newtoggle && !Mgtoggle {
		Mgtoggle = true
		go MultiGameFunc(ctx.Sess, ctx.Conf)
	}

	editConfigfile(ctx.Conf)

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Toggled MultiGame to **%s**", strconv.FormatBool(newtoggle))
	ctx.SendEm(em)
}

func (mgt *MultiGameToggle) Description() string { return "Toggles Multigame" }
func (mgt *MultiGameToggle) Usage() string       { return "" }
func (mgt *MultiGameToggle) Detailed() string {
	return "Toggles if multigame is on or off."
}
func (mgt *MultiGameToggle) Subcommands() map[string]Command { return make(map[string]Command) }

type MultiGame struct{}

func (mg *MultiGame) Message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) == 0 {
		em.Description = "Command `multigame` requires a subcommand!"
	} else {
		em.Description = fmt.Sprintf(`Subcommand "%s" doesn't exist!`, ctx.Args[0])
	}
	ctx.SendEm(em)
}

func (mg *MultiGame) Description() string { return `Commands for Multigame file` }
func (mg *MultiGame) Usage() string       { return "" }
func (mg *MultiGame) Detailed() string {
	return "Commands related to multigame, the timer based 'playing' changer."
}
func (mg *MultiGame) Subcommands() map[string]Command {
	return map[string]Command{"add": &AddMultiGameString{}, "remove": &RemoveMultiGameString{}, "list": &MultiGameList{}, "timer": &MultigameTimer{}, "toggle": &MultiGameToggle{}}
}
