package main

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Message(*Context)
	Description() string
	Usage() string
	Detailed() string
	Subcommands() map[string]Command
}

type CommandHandler struct {
	Commands map[string]Command
}

func (ch *CommandHandler) AddCommand(n string, c Command) {
	ch.Commands[n] = c
}

func HandleSubcommands(ctx *Context, called Command) (*Context, Command) {
	if len(ctx.Args) != 0 {
		scalled, sok := called.Subcommands()[ctx.Args[0]]
		if sok {
			ctx.Invoked += " " + ctx.Args[0]
			ctx.Args = ctx.Args[1:]
			return HandleSubcommands(ctx, scalled)
		} else {
			return ctx, called
		}
	} else {
		return ctx, called
	}
}

func (ch *CommandHandler) HandleCommands(ctx *Context) {
	if ctx.Invoked == "help" {
		ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
		go ch.HelpFunction(ctx)
	} else {
		called, ok := ch.Commands[ctx.Invoked]
		if ok {
			ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
			rctx, rcalled := HandleSubcommands(ctx, called)
			rcalled.Message(rctx)
		} else {
			logerror(errors.New(`Command "` + ctx.Invoked + `" not found`))
		}
	}
}

func (ch *CommandHandler) HelpFunction(ctx *Context) {
	embed := createEmbed(ctx)
	var desc string
	if len(ctx.Args) != 0 {
		ctx.Invoked = ""
		command := ctx.Args[0]
		called, ok := ch.Commands[command]
		ctx.Args = ctx.Args[1:]
		if ok {
			sctx, scalled := HandleSubcommands(ctx, called)
			desc = fmt.Sprintf("`%s%s %s`\n%s", conf.Prefix, command+sctx.Invoked, scalled.Usage(), scalled.Detailed())
			for k, v := range scalled.Subcommands() {
				desc += fmt.Sprintf("\n`%s%s %s` - %s", conf.Prefix, command, k, v.Description())
			}
		} else {
			desc = "No command called `" + command + "` found!"
		}
	} else {
		desc = "Commands:"
		desc += fmt.Sprintf(" `%shelp [command]` for more info!", conf.Prefix)
		for k, v := range ch.Commands {
			desc += fmt.Sprintf("\n`%s%s` - %s", conf.Prefix, k, v.Description())
		}
	}
	embed.Author = &discordgo.MessageEmbedAuthor{Name: ctx.Mess.Author.Username, IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", ctx.Mess.Author.ID, ctx.Mess.Author.Avatar)}
	embed.Description = desc
	embed.Description += "\n\nFloSelfbot [v" + version + "](https://github.com/Moonlington/FloSelfbot)"
	ctx.Sess.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
}
