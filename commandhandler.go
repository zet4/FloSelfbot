package main

import (
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
		ch.Commands[ctx.Invoked].Message(ctx)
	}
}

func (ch *CommandHandler) HelpFunction(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
	color := getUserColor(ctx.Sess, ctx.Guild, ctx.Message.Author.ID)

	var desc string
	desc = "Commands:"

	for k, v := range ch.Commands {
		desc += fmt.Sprintf("\n`%s%s` - %s", conf.Prefix, k, v.Description())
	}

	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: ctx.Message.Author.Username + " - FloSelfbot help", IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", ctx.Message.Author.ID, ctx.Message.Author.Avatar)}, Description: desc, Color: color}
	ctx.Sess.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
}
