package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	AFKMessages        []*discordgo.MessageCreate
	AFKstring          string
	AFKMode            bool
	AFKMultigameBefore bool
)

type AfkPlaying struct{}

func (a *AfkPlaying) Message(ctx *Context) {
	newtoggle := !ctx.Conf.AFKPlay
	ctx.Conf.AFKPlay = newtoggle

	editConfigfile(ctx.Conf)

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Toggled AFKPlay to **%s**", strconv.FormatBool(newtoggle))
	ctx.SendEm(em)
}

func (a *AfkPlaying) Description() string {
	return `Toggles if you want your AFK message to be "played"`
}
func (a *AfkPlaying) Usage() string { return "" }
func (a *AfkPlaying) Detailed() string {
	return `Toggles if you want your AFK message to be "played"`
}
func (a *AfkPlaying) Subcommands() map[string]Command { return make(map[string]Command) }

type Afk struct{}

func (a *Afk) Message(ctx *Context) {
	em := createEmbed(ctx)
	if AFKMode {
		AFKMode = false
		AFKstring = ""
		em.Description = "AFKMode is now off!"
		ctx.Conf.MultigameToggled = AFKMultigameBefore
		var emfields []*discordgo.MessageEmbedField
		for _, msg := range AFKMessages {
			field := &discordgo.MessageEmbedField{Inline: false, Name: msg.Author.Username + " in <#" + msg.ChannelID + ">", Value: msg.Content}
			emfields = append(emfields, field)
		}
		em.Fields = emfields
		ctx.SendEm(em)
		AFKMessages = []*discordgo.MessageCreate{}
	} else {
		AFKMode = true
		AFKstring = strings.Join(ctx.Args, " ")
		if ctx.Conf.AFKPlay {
			var txt string
			if AFKstring != "" {
				txt = "AFK - " + AFKstring
			} else {
				txt = "AFK"
			}
			currentgame = txt
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

func (a *Afk) Description() string { return `Sets your selfbot to "AFK Mode"` }
func (a *Afk) Usage() string       { return "[message]" }
func (a *Afk) Detailed() string {
	return "Lets people know when you are AFK (Might be removed soon cuz discord selfbot guidelines)"
}
func (a *Afk) Subcommands() map[string]Command { return map[string]Command{"playing": &AfkPlaying{}} }
