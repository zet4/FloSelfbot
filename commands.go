package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/robertkrimen/otto"
)

func createEmbed(ctx *Context) *discordgo.MessageEmbed {
	if conf.EmbedColor == "#000000" || conf.EmbedColor == "" {
		if conf.EmbedColor == "" {
			conf.EmbedColor = "#000000"
			editConfigfile(conf)
		}
		color := ctx.Sess.State.UserColor(ctx.Mess.Author.ID, ctx.Mess.ChannelID)
		return &discordgo.MessageEmbed{Color: color}
	} else {
		color := conf.EmbedColor
		if strings.HasPrefix(color, "#") {
			color = "0x" + conf.EmbedColor[1:]
		}
		d, _ := strconv.ParseInt(color, 0, 64)
		return &discordgo.MessageEmbed{Color: int(d)}
	}
}

type Ping struct{}

func (p *Ping) Message(ctx *Context) {
	em := createEmbed(ctx)
	em.Description = "Pong!"
	start := time.Now()
	msg, _ := ctx.SendEm(em)
	elapsed := time.Since(start)
	em.Description = fmt.Sprintf("Pong! `%s`", elapsed)
	ctx.Sess.ChannelMessageEditEmbed(ctx.Mess.ChannelID, msg.ID, em)
}

func (p *Ping) Description() string             { return "Measures latency" }
func (p *Ping) Usage() string                   { return "" }
func (p *Ping) Detailed() string                { return "Measures latency" }
func (p *Ping) Subcommands() map[string]Command { return make(map[string]Command) }

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

type Me struct{}

func (m *Me) Message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) != 0 {
		text := strings.Join(ctx.Args, " ")
		em.Description = fmt.Sprintf("***%s*** *%s*", ctx.Mess.Author.Username, text)
		ctx.SendEm(em)
	} else {
		em.Description = fmt.Sprintf("***%s*** *was silent...*", ctx.Mess.Author.Username)
		ctx.SendEm(em)
	}
}

func (m *Me) Description() string             { return "Says stuff" }
func (m *Me) Usage() string                   { return "<message>" }
func (m *Me) Detailed() string                { return "Says stuff." }
func (m *Me) Subcommands() map[string]Command { return make(map[string]Command) }

type Embed struct{}

func (e *Embed) Message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) != 0 {
		text := strings.Join(ctx.Args, " ")
		em.Description = fmt.Sprintf("%s", text)
		ctx.SendEm(em)
	} else {
		em.Description = fmt.Sprintf("***%s*** *was silent...*", ctx.Mess.Author.Username)
		ctx.SendEm(em)
	}
}

func (e *Embed) Description() string             { return "Embeds stuff" }
func (e *Embed) Usage() string                   { return "<message>" }
func (e *Embed) Detailed() string                { return "Embeds stuff." }
func (e *Embed) Subcommands() map[string]Command { return make(map[string]Command) }

type Eval struct{}

func (e *Eval) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		vm := otto.New()
		vm.Set("ctx", ctx)
		vm.Set("conf", conf)
		toEval := strings.Join(ctx.Args, " ")
		executed, err := vm.Run(toEval)
		em := createEmbed(ctx)
		if err != nil {
			em.Description = fmt.Sprintf("Input: `%s`\n\nError: `%s`", toEval, err.Error())
			ctx.SendEm(em)
			return
		}
		em.Description = fmt.Sprintf("Input: `%s`\n\nOutput: ```js\n%s\n```", toEval, executed.String())
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't eval shit"
		ctx.SendEm(em)
	}
}

func (e *Eval) Description() string { return "Evaluates using Otto (Advanced stuff, don't bother)" }
func (e *Eval) Usage() string       { return "<toEval>" }
func (e *Eval) Detailed() string {
	return "I'm serious, don't bother with this command."
}
func (e *Eval) Subcommands() map[string]Command { return make(map[string]Command) }

type Clean struct{}

func (c *Clean) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
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
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify an amount"
		ctx.SendEm(em)
	}
}

func (c *Clean) Description() string { return "Cleans up your messages" }
func (c *Clean) Usage() string       { return "<amount>" }
func (c *Clean) Detailed() string {
	return "If you realise you have been spamming a little, this is the command to use then."
}
func (c *Clean) Subcommands() map[string]Command { return make(map[string]Command) }

type Quote struct{}

func (q *Quote) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		var qmess *discordgo.Message

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
		timestampo := timestamp.Format(time.ANSIC)
		emfooter := &discordgo.MessageEmbedFooter{Text: "Sent | " + timestampo}
		emcolor := ctx.Sess.State.UserColor(qmess.Author.ID, qmess.ChannelID)
		em := &discordgo.MessageEmbed{Author: emauthor, Footer: emfooter, Description: qmess.Content, Color: emcolor}
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't specify a message ID"
		ctx.SendEm(em)
	}
}

func (q *Quote) Description() string { return "Quotes a message from the last 100 messages" }
func (q *Quote) Usage() string       { return "<messageID>" }
func (q *Quote) Detailed() string {
	return "To find messageID you first need to turn on Developer mode in discord, then right click any message and click 'Copy ID'"
}
func (q *Quote) Subcommands() map[string]Command { return make(map[string]Command) }

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

type ChangePrefix struct{}

func (cp *ChangePrefix) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		newprefix := strings.Join(ctx.Args, " ")
		conf.Prefix = newprefix

		editConfigfile(conf)

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
			conf.EmbedColor = ctx.Args[0]
			editConfigfile(conf)
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
	newlogmode := !conf.LogMode
	conf.LogMode = newlogmode

	editConfigfile(conf)

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
	toml.DecodeFile("config.toml", &conf)
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

type AddMultiGameString struct{}

func (a *AddMultiGameString) Message(ctx *Context) {
	if len(ctx.Args) != 0 {
		em := createEmbed(ctx)
		game := strings.Join(ctx.Args, " ")
		conf.MultiGameStrings = append(conf.MultiGameStrings, game)
		editConfigfile(conf)
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
		for i, v := range conf.MultiGameStrings {
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
			conf.MultiGameStrings = append(conf.MultiGameStrings[:pos], conf.MultiGameStrings[pos+1:]...)
			editConfigfile(conf)
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
	for _, v := range conf.MultiGameStrings {
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
		conf.MultiGameMinutes = delay
		editConfigfile(conf)
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
	newtoggle := !conf.MultigameToggled
	conf.MultigameToggled = newtoggle

	if newtoggle && !Mgtoggle {
		Mgtoggle = true
		go MultiGameFunc(ctx.Sess)
	}

	editConfigfile(conf)

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

// type AddEmote struct{}
//
// func (a *AddEmote) Message(ctx *Context) {
// 	if len(ctx.Args) <= 2 {
// 		em := createEmbed(ctx)
// 		emote := strings.Join(ctx.Args[1:], " ")
// 		conf.Emotes[ctx.Args[0]] = emote
// 		editConfigfile(conf)
// 		em.Description = fmt.Sprintf("Added **%s -> %s** to Emotes", ctx.Args[0], emote)
// 		ctx.SendEm(em)
// 	} else {
// 		em := createEmbed(ctx)
// 		em.Description = "No Emote name and/or content specified!"
// 		ctx.SendEm(em)
// 	}
// }
//
// func (a *AddEmote) Description() string { return "Adds a string to Emote" }
// func (a *AddEmote) Usage() string       { return "<game>" }
// func (a *AddEmote) Detailed() string {
// 	return "Adds a string to Emote."
// }
// func (a *AddEmote) Subcommands() map[string]Command { return make(map[string]Command) }
//
// type RemoveEmote struct{}
//
// func (r *RemoveEmote) Message(ctx *Context) {
// 	if len(ctx.Args) != 0 {
// 		delete(conf.Emotes, ctx.Args[0])
// 	} else {
// 		em := createEmbed(ctx)
// 		em.Description = "No Emote name specified!"
// 		ctx.SendEm(em)
// 	}
// }
//
// func (r *RemoveEmote) Description() string { return "Removes a string from Emote" }
// func (r *RemoveEmote) Usage() string       { return "<game>" }
// func (r *RemoveEmote) Detailed() string {
// 	return "Removes a string from Emote."
// }
// func (a *RemoveEmote) Subcommands() map[string]Command { return make(map[string]Command) }
//
// type EmoteList struct{}
//
// func (l *EmoteList) Message(ctx *Context) {
// 	em := createEmbed(ctx)
// 	desc := "Current strings in Emote:"
// 	for k, v := range conf.Emotes {
// 		desc += fmt.Sprintf("\n`%s` -> `%s`", k, v)
// 	}
// 	em.Description = desc
// 	ctx.SendEm(em)
// }
//
// func (l *EmoteList) Description() string { return "Returns your Emote strings" }
// func (l *EmoteList) Usage() string       { return "" }
// func (l *EmoteList) Detailed() string {
// 	return "Returns the list of all your strings in Emote"
// }
// func (l *EmoteList) Subcommands() map[string]Command { return make(map[string]Command) }
//
// type Emote struct{}
//
// func (e *Emote) Message(ctx *Context) {
// 	em := createEmbed(ctx)
// 	if len(ctx.Args) == 0 {
// 		em.Description = "Command `emote` requires a subcommand!"
// 	} else {
// 		em.Description = fmt.Sprintf(`Subcommand "%s" doesn't exist!`, ctx.Args[0])
// 	}
// 	ctx.SendEm(em)
// }
//
// func (e *Emote) Description() string { return `Commmands related to emotes` }
// func (e *Emote) Usage() string       { return "" }
// func (e *Emote) Detailed() string {
// 	return "Commmands related to emotes, trigger them by doing :nameofemote:"
// }
// func (e *Emote) Subcommands() map[string]Command {
// 	return map[string]Command{"add": &AddEmote{}, "remove": &RemoveEmote{}, "list": &EmoteList{}}
// }
