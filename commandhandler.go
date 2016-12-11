package main

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Message(*Context)
	Description() string
}

type CommandHandler struct {
	Commands map[string]Command
}

func (ch *CommandHandler) AddCommand(n string, c Command) {
	ch.Commands[n] = c
}

func (ch *CommandHandler) HandleCommands(ctx *Context) {
	if ctx.Invoked == "help" {
		ch.HelpFunction(ctx)
	} else {
		called, ok := ch.Commands[ctx.Invoked]
		if ok {
			called.Message(ctx)
		} else {
			logerror(errors.New(`Command "` + ctx.Invoked + `" not found`))
		}
	}
}

func (ch *CommandHandler) HelpFunction(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	color := ctx.Sess.State.UserColor(ctx.Mess.Author.ID, ctx.Mess.ChannelID)

	var desc string
	desc = "Commands:"

	for k, v := range ch.Commands {
		desc += fmt.Sprintf("\n`%s%s` - %s", conf.Prefix, k, v.Description())
	}

	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: ctx.Mess.Author.Username + " - FloSelfbot help", IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", ctx.Mess.Author.ID, ctx.Mess.Author.Avatar)}, Description: desc, Color: color}
	ctx.Sess.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
}
