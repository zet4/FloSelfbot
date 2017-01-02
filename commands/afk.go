package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	AFKMessages []*discordgo.MessageCreate
	AFKstring   string
	AFKMode     bool
)

type Afk struct{}

func (a *Afk) Message(ctx *Context) {
	em := createEmbed(ctx)
	if AFKMode {
		AFKMode = false
		AFKstring = ""
		em.Description = "AFKMode is now off!"
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
		em.Description = "AFKMode is now on!"
		ctx.SendEm(em)
	}
}

func (a *Afk) Description() string { return `Sets your selfbot to "AFK Mode"` }
func (a *Afk) Usage() string       { return "[message]" }
func (a *Afk) Detailed() string {
	return "Lets people know when you are AFK (Might be removed soon cuz discord selfbot guidelines)"
}
func (a *Afk) Subcommands() map[string]Command { return make(map[string]Command) }
