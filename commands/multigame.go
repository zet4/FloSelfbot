package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Mgtoggle bool that handles if Multigame is Toggled
var Mgtoggle bool

// MultiGameFunc handles Multigame
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
			EditConfigFile(conf)
		}
		time.Sleep(time.Minute * time.Duration(conf.MultiGameMinutes))
	}
}

type addMultiGameString struct{}

func (a *addMultiGameString) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		game := strings.Join(ctx.Args, " ")
		if strings.ToLower(game) == "all" {
			em.Description = fmt.Sprintf("**All** can not be added because of problems with the remove command.", game)
			ctx.SendEm(em)
			return
		}
		ctx.Conf.MultiGameStrings = append(ctx.Conf.MultiGameStrings, game)
		EditConfigFile(ctx.Conf)
		em.Description = fmt.Sprintf("Added **%s** to Multigame", game)
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "No game name specified!"
		ctx.SendEm(em)
	}
}

func (a *addMultiGameString) description() string { return "Adds a string to the Multigame" }
func (a *addMultiGameString) usage() string       { return "<game>" }
func (a *addMultiGameString) detailed() string {
	return "Adds a string to the Multigame."
}
func (a *addMultiGameString) subcommands() map[string]Command { return make(map[string]Command) }

type removeMultiGameString struct{}

func (r *removeMultiGameString) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		var pos int
		em := createEmbed(ctx)
		game := strings.Join(ctx.Args, " ")
		if strings.ToLower(game) == "all" {
			ctx.Conf.MultiGameStrings = []string{}
			EditConfigFile(ctx.Conf)
			em.Description = "Removed everything from the list."
			ctx.SendEm(em)
			return
		}
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
			EditConfigFile(ctx.Conf)
			em.Description = fmt.Sprintf("Removed **%s** from Multigame", game)
			ctx.SendEm(em)
		}
	} else {
		em := createEmbed(ctx)
		em.Description = "No game name specified!"
		ctx.SendEm(em)
	}
}

func (r *removeMultiGameString) description() string { return "Removes a string from Multigame" }
func (r *removeMultiGameString) usage() string       { return "<game|all>" }
func (r *removeMultiGameString) detailed() string {
	return "Removes a string from Multigame. Type All as game to remove everything from the list"
}
func (r *removeMultiGameString) subcommands() map[string]Command { return make(map[string]Command) }

type multiGameList struct{}

func (l *multiGameList) message(ctx *Context) {
	var desc string
	em := createEmbed(ctx)
	if len(ctx.Conf.MultiGameStrings) == 0 {
		desc = "There are currently no games in Multigame!"
	} else {
		desc = "Current strings in Multigame:"
		for _, v := range ctx.Conf.MultiGameStrings {
			desc += fmt.Sprintf("\n`%s`", v)
		}
	}
	desc += fmt.Sprintf("\nCurrent timer is set to: **%s**", strconv.Itoa(ctx.Conf.MultiGameMinutes))
	em.Description = desc
	ctx.SendEm(em)
}

func (l *multiGameList) description() string { return "Returns your Multigame strings" }
func (l *multiGameList) usage() string       { return "" }
func (l *multiGameList) detailed() string {
	return "Returns the list of all your strings in Multigame"
}
func (l *multiGameList) subcommands() map[string]Command { return make(map[string]Command) }

type multigameTimer struct{}

func (mgt *multigameTimer) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		delay, err := strconv.Atoi(ctx.Args[0])
		if err != nil {
			em.Description = "Invalid time!"
			ctx.SendEm(em)
			return
		}
		ctx.Conf.MultiGameMinutes = delay
		EditConfigFile(ctx.Conf)
		em.Description = fmt.Sprintf("Multigame timer set to **%s** minutes", ctx.Args[0])
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "No int specified!"
		ctx.SendEm(em)
	}
}

func (mgt *multigameTimer) description() string {
	return "Changes the delay between switches (in minutes)"
}
func (mgt *multigameTimer) usage() string { return "<minutes>" }
func (mgt *multigameTimer) detailed() string {
	return "Changes the delay between switches (in minutes)"
}
func (mgt *multigameTimer) subcommands() map[string]Command { return make(map[string]Command) }

type multiGameToggle struct{}

func (mgt *multiGameToggle) message(ctx *Context) {
	newtoggle := !ctx.Conf.MultigameToggled
	ctx.Conf.MultigameToggled = newtoggle

	if newtoggle && !Mgtoggle {
		Mgtoggle = true
		go MultiGameFunc(ctx.Sess, ctx.Conf)
	}

	EditConfigFile(ctx.Conf)

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Toggled MultiGame to **%s**", strconv.FormatBool(newtoggle))
	ctx.SendEm(em)
}

func (mgt *multiGameToggle) description() string { return "Toggles Multigame" }
func (mgt *multiGameToggle) usage() string       { return "" }
func (mgt *multiGameToggle) detailed() string {
	return "Toggles if multigame is on or off."
}
func (mgt *multiGameToggle) subcommands() map[string]Command { return make(map[string]Command) }

// MultiGame struct handles MultiGame Command
type MultiGame struct{}

func (mg *MultiGame) message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) == 0 {
		em.Description = "Command `multigame` requires a subcommand!"
	} else {
		em.Description = fmt.Sprintf(`Subcommand "%s" doesn't exist!`, ctx.Args[0])
	}
	ctx.SendEm(em)
}

func (mg *MultiGame) description() string { return `Commands for Multigame file` }
func (mg *MultiGame) usage() string       { return "" }
func (mg *MultiGame) detailed() string {
	return "Commands related to multigame, the timer based 'playing' changer."
}
func (mg *MultiGame) subcommands() map[string]Command {
	return map[string]Command{"add": &addMultiGameString{}, "remove": &removeMultiGameString{}, "list": &multiGameList{}, "timer": &multigameTimer{}, "toggle": &multiGameToggle{}}
}
