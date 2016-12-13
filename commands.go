package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/robertkrimen/otto"
)

func createEmbed(ctx *Context) *discordgo.MessageEmbed {
	color := ctx.Sess.State.UserColor(ctx.Mess.Author.ID, ctx.Mess.ChannelID)
	return &discordgo.MessageEmbed{Color: color}
}

// func (s *Name) Message(ctx *Context) {
//
// }
//
// func (s *Name) Description() string { return "" }
// func (s *Name) Usage() string       { return "" }
// func (s *Name) Detailed() string    { return "" }

type Ping struct{}

func (p *Ping) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	em := createEmbed(ctx)
	em.Description = "Pong!"
	start := time.Now()
	msg, _ := ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
	elapsed := time.Since(start)
	em.Description = fmt.Sprintf("Pong! `%s`", elapsed)
	ctx.Sess.ChannelMessageEditEmbed(ctx.Mess.ChannelID, msg.ID, em)
}

func (p *Ping) Description() string { return "Measures latency" }
func (p *Ping) Usage() string       { return "" }
func (p *Ping) Detailed() string    { return "Measures latency" }

type SetGame struct{}

func (sg *SetGame) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	em := createEmbed(ctx)
	game := strings.Join(ctx.Args, " ")
	em.Description = fmt.Sprintf("Changed game to **%s**", game)
	ctx.Sess.UpdateStatus(0, game)
	ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)

}

func (sg *SetGame) Description() string { return "Sets your game to anything you like" }
func (sg *SetGame) Usage() string       { return "<game>" }
func (sg *SetGame) Detailed() string {
	return "Changes your 'Playing' status on discord (Because of discord you cant see the change yourself.)"
}

type Me struct{}

func (m *Me) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	em := createEmbed(ctx)
	text := strings.Join(ctx.Args, " ")
	em.Description = fmt.Sprintf("***%s*** *%s*", ctx.Mess.Author.Username, text)
	ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func (m *Me) Description() string { return "Says stuff" }
func (m *Me) Usage() string       { return "<message>" }
func (m *Me) Detailed() string    { return "Says stuff." }

type Eval struct{}

func (e *Eval) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	vm := otto.New()
	vm.Set("ctx", ctx)
	toEval := strings.Join(ctx.Args, " ")
	executed, err := vm.Run(toEval)
	em := createEmbed(ctx)
	if err != nil {
		em.Description = fmt.Sprintf("Input: `%s`\n\nError: `%s`", toEval, err.Error())
		ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
		return
	}
	em.Description = fmt.Sprintf("Input: `%s`\n\nOutput: ```js\n%s\n```", toEval, executed.String())
	ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func (e *Eval) Description() string { return "Evaluates using Otto (Advanced stuff, don't bother)" }
func (e *Eval) Usage() string       { return "<toEval>" }
func (e *Eval) Detailed() string {
	return "I'm serious, don't bother with this command."
}

type Clean struct{}

func (c *Clean) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	limit, err := strconv.Atoi(ctx.Args[0])
	logerror(err)
	msgs, err := ctx.Sess.ChannelMessages(ctx.Mess.ChannelID, limit, ctx.Mess.ID, "")
	logerror(err)
	for _, msg := range msgs {
		if msg.Author.ID == ctx.Sess.State.User.ID {
			err = ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, msg.ID)
			logerror(err)
		}
	}
}

func (c *Clean) Description() string { return "Cleans up your messages" }
func (c *Clean) Usage() string       { return "<amount>" }
func (c *Clean) Detailed() string {
	return "If you realise you have been spamming a little, this is the command to use then."
}

type Quote struct{}

func (q *Quote) Message(ctx *Context) {
	var qmess *discordgo.Message
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	mID := ctx.Args[0]
	msgs, err := ctx.Sess.ChannelMessages(ctx.Mess.ChannelID, 100, ctx.Mess.ID, "")
	for _, msg := range msgs {
		if msg.ID == mID {
			qmess = msg
		}
	}
	if qmess == nil {
		ctx.Sess.ChannelMessageSend(ctx.Mess.ChannelID, "Message not found in last 100 messages.")
	}

	emauthor := &discordgo.MessageEmbedAuthor{Name: qmess.Author.Username, IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", qmess.Author.ID, qmess.Author.Avatar)}
	timestamp, err := qmess.Timestamp.Parse()
	logerror(err)
	emfooter := &discordgo.MessageEmbedFooter{Text: "Sent | " + timestamp.String()}
	emcolor := ctx.Sess.State.UserColor(qmess.Author.ID, qmess.ChannelID)
	em := &discordgo.MessageEmbed{Author: emauthor, Footer: emfooter, Description: qmess.Content, Color: emcolor}
	ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func (q *Quote) Description() string { return "Quotes a message from the last 100 messages" }
func (q *Quote) Usage() string       { return "<messageID>" }
func (q *Quote) Detailed() string {
	return "To find messageID you first need to turn on Developer mode in discord, then right click any message and click 'Copy ID'"
}

type Afk struct{}

func (a *Afk) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
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
		ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
		AFKMessages = []*discordgo.MessageCreate{}
	} else {
		AFKMode = true
		AFKstring = strings.Join(ctx.Args, " ")
		em.Description = "AFKMode is now on!"
		ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
	}
}

func (a *Afk) Description() string { return `Sets your selfbot to "AFK Mode"` }
func (a *Afk) Usage() string       { return "[message]" }
func (a *Afk) Detailed() string {
	return "Lets people know when you are AFK (Might be removed soon cuz discord selfbot guidelines)"
}

type ChangePrefix struct{}

func (cp *ChangePrefix) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	newprefix := strings.Join(ctx.Args, " ")
	conf.Prefix = newprefix

	f, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	logwarning(err)
	defer f.Close()
	logwarning(toml.NewEncoder(f).Encode(conf))

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Changed prefix to **%s**", newprefix)
	ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func (cp *ChangePrefix) Description() string { return `Changes your prefix` }
func (cp *ChangePrefix) Usage() string       { return "<newprefix>" }
func (cp *ChangePrefix) Detailed() string {
	return "Changes your prefix (You can do the same by editing the config.toml file)"
}

type ToggleLogMode struct{}

func (l *ToggleLogMode) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	newlogmode := !conf.LogMode
	conf.LogMode = newlogmode

	f, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	logwarning(err)
	defer f.Close()
	logwarning(toml.NewEncoder(f).Encode(conf))

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Toggled LogMode to **%s**", strconv.FormatBool(newlogmode))
	ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func (l *ToggleLogMode) Description() string { return `Toggles logmode` }
func (l *ToggleLogMode) Usage() string       { return "" }
func (l *ToggleLogMode) Detailed() string {
	return "Toggles Logmode on or off (You can do the same by editing the config.toml file)"
}
