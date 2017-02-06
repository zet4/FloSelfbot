package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type changePrefix struct{}

func (cp *changePrefix) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		newprefix := strings.Join(ctx.Args, " ")
		ctx.Conf.Prefix = newprefix

		EditConfigFile(ctx.Conf)

		em := createEmbed(ctx)
		em.Description = fmt.Sprintf("Changed prefix to **%s**", newprefix)
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify a new prefix."
		ctx.SendEm(em)
	}
}

func (cp *changePrefix) description() string { return `Changes your prefix` }
func (cp *changePrefix) usage() string       { return "<newprefix>" }
func (cp *changePrefix) detailed() string {
	return "Changes your prefix (You can do the same by editing the config.toml file)"
}
func (cp *changePrefix) subcommands() map[string]Command { return make(map[string]Command) }

type changeEmbedColor struct{}

func (cec *changeEmbedColor) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		if strings.HasPrefix(ctx.Args[0], "#") {
			ctx.Conf.EmbedColor = ctx.Args[0]
			EditConfigFile(ctx.Conf)
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

func (cec *changeEmbedColor) description() string { return `Changes your embed color` }
func (cec *changeEmbedColor) usage() string       { return "<newhex>" }
func (cec *changeEmbedColor) detailed() string {
	return "Changes your EmbedColor (You can do the same by editing the config.toml file)"
}
func (cec *changeEmbedColor) subcommands() map[string]Command { return make(map[string]Command) }

type toggleLogMode struct{}

func (l *toggleLogMode) message(ctx *Context) {
	newlogmode := !ctx.Conf.LogMode
	ctx.Conf.LogMode = newlogmode

	EditConfigFile(ctx.Conf)

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Toggled LogMode to **%s**", strconv.FormatBool(newlogmode))
	ctx.SendEm(em)
}

func (l *toggleLogMode) description() string { return `Toggles logmode` }
func (l *toggleLogMode) usage() string       { return "" }
func (l *toggleLogMode) detailed() string {
	return "Toggles Logmode on or off (You can do the same by editing the config.toml file)"
}
func (l *toggleLogMode) subcommands() map[string]Command { return make(map[string]Command) }

type reloadConfig struct{}

func (r *reloadConfig) message(ctx *Context) {
	toml.DecodeFile("config.toml", &ctx.Conf)
	desc := "Reloaded Config file!\n"
	em := createEmbed(ctx)
	em.Description = desc
	ctx.SendEm(em)
}

func (r *reloadConfig) description() string { return "Reloads the config" }
func (r *reloadConfig) usage() string       { return "" }
func (r *reloadConfig) detailed() string {
	return "If you made any changes to the config file you dont have to restart the selfbot."
}
func (r *reloadConfig) subcommands() map[string]Command { return make(map[string]Command) }

type changeAutoDeleteTimer struct{}

func (cadt *changeAutoDeleteTimer) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		delay, err := strconv.Atoi(ctx.Args[0])
		if err != nil {
			em.Description = "Invalid time!"
			ctx.SendEm(em)
			return
		}
		ctx.Conf.AutoDeleteSeconds = delay
		EditConfigFile(ctx.Conf)
		em.Description = fmt.Sprintf("Autodelete timer set to **%s** seconds", ctx.Args[0])
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "No int specified!"
		ctx.SendEm(em)
	}
}
func (cadt *changeAutoDeleteTimer) description() string {
	return "Changes autodelete timer (in seconds)"
}
func (cadt *changeAutoDeleteTimer) usage() string { return "<seconds>" }
func (cadt *changeAutoDeleteTimer) detailed() string {
	return "Changes how many seconds it takes for the selfbot to delete its post (in seconds) For it to not delete, use 0."
}
func (cadt *changeAutoDeleteTimer) subcommands() map[string]Command { return make(map[string]Command) }

// Configcommand struct handles Config Command
type Configcommand struct{}

func (c *Configcommand) message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) == 0 {
		em.Description = "Command `config` requires a subcommand!"
	} else {
		em.Description = fmt.Sprintf(`Subcommand "%s" doesn't exist!`, ctx.Args[0])
	}
	ctx.SendEm(em)
}

func (c *Configcommand) description() string { return `Commands for the config file` }
func (c *Configcommand) usage() string       { return "" }
func (c *Configcommand) detailed() string {
	return "Commands related to the config file, like changing color, changing prefix, etc."
}
func (c *Configcommand) subcommands() map[string]Command {
	return map[string]Command{"togglelogmode": &toggleLogMode{}, "reload": &reloadConfig{}, "changeprefix": &changePrefix{}, "changecolor": &changeEmbedColor{}, "setautodelete": &changeAutoDeleteTimer{}}
}
