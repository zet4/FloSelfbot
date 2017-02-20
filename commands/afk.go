package commands

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	// AFKMessages contains messages in AFK
	AFKMessages []*discordgo.MessageCreate
	// AFKstring user-set string that contains why the user is AFK
	AFKstring string
	// AFKMode bool that says if AFKMode is on
	AFKMode bool
	// AFKMultigameBefore handles if Multigame was on before AFK
	AFKMultigameBefore bool
)

type afkPlaying struct{}

func (a *afkPlaying) message(ctx *Context) {
	newtoggle := !ctx.Conf.AFKPlay
	ctx.Conf.AFKPlay = newtoggle

	EditConfigFile(ctx.Conf)

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Toggled AFKPlay to **%s**", strconv.FormatBool(newtoggle))
	ctx.SendEm(em)
}

func (a *afkPlaying) description() string {
	return `Toggles if you want your AFK message to be "played"`
}
func (a *afkPlaying) usage() string { return "" }
func (a *afkPlaying) detailed() string {
	return `Toggles if you want your AFK message to be "played"`
}
func (a *afkPlaying) subcommands() map[string]Command { return make(map[string]Command) }

// Afk struct handles Afk Command
type Afk struct{}

func (a *Afk) message(ctx *Context) {
	em := createEmbed(ctx)
	if AFKMode {
		AFKMode = false
		AFKstring = ""
		em.Description = "AFKMode is now off!"
		ctx.Conf.MultigameToggled = AFKMultigameBefore
		if ctx.Conf.AFKPlay {
			ctx.Sess.UpdateStatus(0, "")
			CurrentGame = ""
		}
		var emfields []*discordgo.MessageEmbedField
		for _, msg := range AFKMessages {
			field := &discordgo.MessageEmbedField{Inline: false, Name: msg.Author.Username + " in <#" + msg.ChannelID + ">", Value: msg.Content}
			emfields = append(emfields, field)
		}
		em.Fields = emfields
		ctx.SendEm(em)
		AFKMessages = []*discordgo.MessageCreate{}
		ctx.Sess.UserUpdateStatus(discordgo.StatusOnline)
	} else {
		AFKMode = true
		AFKstring = ctx.Argstr
		if ctx.Conf.AFKPlay {
			var txt string
			if AFKstring != "" {
				txt = "AFK - " + AFKstring
			} else {
				txt = "AFK"
			}
			CurrentGame = txt
			err := ctx.Sess.UpdateStatus(0, txt)
			logerror(err)
			AFKMultigameBefore = ctx.Conf.MultigameToggled
			ctx.Conf.MultigameToggled = false
		}
		em.Description = "AFKMode is now on!"
		ctx.SendEm(em)
		ctx.Sess.UserUpdateStatus(discordgo.StatusDoNotDisturb)
	}
}

func (a *Afk) description() string { return `Sets your selfbot to "AFK Mode"` }
func (a *Afk) usage() string       { return "[message]" }
func (a *Afk) detailed() string {
	return "Lets people know when you are AFK (Might be removed soon cuz discord selfbot guidelines)"
}
func (a *Afk) subcommands() map[string]Command { return map[string]Command{"playing": &afkPlaying{}} }
