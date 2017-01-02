package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type ChangePrefix struct{}

func (cp *ChangePrefix) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		newprefix := strings.Join(ctx.Args, " ")
		ctx.Conf.Prefix = newprefix

		editConfigfile(ctx.Conf)

		em := createEmbed(ctx)
		em.Description = fmt.Sprintf("Changed prefix to **%s**", newprefix)
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify a new prefix."
		ctx.SendEm(em)
	}
}

func (cp *ChangePrefix) Description() string { return `Changes your prefix` }
func (cp *ChangePrefix) Usage() string       { return "<newprefix>" }
func (cp *ChangePrefix) Detailed() string {
	return "Changes your prefix (You can do the same by editing the config.toml file)"
}
func (cp *ChangePrefix) Subcommands() map[string]Command { return make(map[string]Command) }

type ChangeEmbedColor struct{}

func (cec *ChangeEmbedColor) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		if strings.HasPrefix(ctx.Args[0], "#") {
			ctx.Conf.EmbedColor = ctx.Args[0]
			editConfigfile(ctx.Conf)
			em := createEmbed(ctx)
			em.Description = fmt.Sprintf("Changed EmbedColor to **%s**", ctx.Args[0])
			ctx.SendEm(em)
		} else {
			em := createEmbed(ctx)
			em.Description = "Color needs to be in hexadecimal format. `#FA6409` for example. (#000000 for default)"
			ctx.SendEm(em)
		}
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify a new hex value. (#000000 for default)"
		ctx.SendEm(em)
	}
}

func (cec *ChangeEmbedColor) Description() string { return `Changes your embed color` }
func (cec *ChangeEmbedColor) Usage() string       { return "<newhex>" }
func (cec *ChangeEmbedColor) Detailed() string {
	return "Changes your EmbedColor (You can do the same by editing the config.toml file)"
}
func (cec *ChangeEmbedColor) Subcommands() map[string]Command { return make(map[string]Command) }

type ToggleLogMode struct{}

func (l *ToggleLogMode) Message(ctx *Context) {
	newlogmode := !ctx.Conf.LogMode
	ctx.Conf.LogMode = newlogmode

	editConfigfile(ctx.Conf)

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Toggled LogMode to **%s**", strconv.FormatBool(newlogmode))
	ctx.SendEm(em)
}

func (l *ToggleLogMode) Description() string { return `Toggles logmode` }
func (l *ToggleLogMode) Usage() string       { return "" }
func (l *ToggleLogMode) Detailed() string {
	return "Toggles Logmode on or off (You can do the same by editing the config.toml file)"
}
func (l *ToggleLogMode) Subcommands() map[string]Command { return make(map[string]Command) }

type ReloadConfig struct{}

func (r *ReloadConfig) Message(ctx *Context) {
	toml.DecodeFile("config.toml", &ctx.Conf)
	desc := "Reloaded Config file!\n"
	em := createEmbed(ctx)
	em.Description = desc
	ctx.SendEm(em)
}

func (r *ReloadConfig) Description() string { return "Reloads the config" }
func (r *ReloadConfig) Usage() string       { return "" }
func (r *ReloadConfig) Detailed() string {
	return "If you made any changes to the config file you dont have to restart the selfbot."
}
func (r *ReloadConfig) Subcommands() map[string]Command { return make(map[string]Command) }

type Configcommand struct{}

func (c *Configcommand) Message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) == 0 {
		em.Description = "Command `config` requires a subcommand!"
	} else {
		em.Description = fmt.Sprintf(`Subcommand "%s" doesn't exist!`, ctx.Args[0])
	}
	ctx.SendEm(em)
}

func (c *Configcommand) Description() string { return `Commands for the config file` }
func (c *Configcommand) Usage() string       { return "" }
func (c *Configcommand) Detailed() string {
	return "Commands related to the config file, like changing color, changing prefix, etc."
}
func (c *Configcommand) Subcommands() map[string]Command {
	return map[string]Command{"togglelogmode": &ToggleLogMode{}, "reload": &ReloadConfig{}, "changeprefix": &ChangePrefix{}, "changecolor": &ChangeEmbedColor{}}
}
