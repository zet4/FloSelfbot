package commands

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// A Command interface stores functions for commands
type Command interface {
	message(*Context)
	description() string
	usage() string
	detailed() string
	subcommands() map[string]Command
}

// A CommandHandler handles Commands
type CommandHandler struct {
	Commands   map[string]Command
	Categories map[string]map[string]Command
}

// AddCommand adds a Command to the CommandHandler
// n: Name for the Command
// c: Command to add
func (ch *CommandHandler) AddCommand(n, category string, c Command) {
	ch.Commands[n] = c
	ch.setCategory(category, n, c)
}

func (ch *CommandHandler) setCategory(category, n string, c Command) {
	child, ok := ch.Categories[category]
	if !ok {
		child = map[string]Command{}
		ch.Categories[category] = child
	}
	child[n] = c
}

// HandleSubcommands returns the Context and Command that is being called
// ctx: Context used
// called: Command called
func HandleSubcommands(ctx *Context, called Command) (*Context, Command) {
	if len(ctx.Args) != 0 {
		scalled, sok := called.subcommands()[strings.ToLower(ctx.Args[0])]
		if sok {
			ctx.Argstr = ctx.Argstr[len(ctx.Args[0]):]
			if ctx.Argstr != "" {
				ctx.Argstr = ctx.Argstr[1:]
			}
			ctx.Invoked += " " + ctx.Args[0]
			ctx.Args = ctx.Args[1:]
			return HandleSubcommands(ctx, scalled)
		}
	}
	return ctx, called
}

// HandleCommands handles the Context and calls Command
// ctx: Context used
func (ch *CommandHandler) HandleCommands(ctx *Context) {
	if strings.ToLower(ctx.Invoked) == "help" {
		ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
		go ch.HelpFunction(ctx)
	} else {
		called, ok := ch.Commands[strings.ToLower(ctx.Invoked)]
		if ok {
			ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
			rctx, rcalled := HandleSubcommands(ctx, called)
			go rcalled.message(rctx)
		} else {
			logerror(errors.New(`Command "` + ctx.Invoked + `" not found`))
		}
	}
}

// HelpFunction handles the Help command for the CommandHandler
// ctx: Context used
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
			desc = fmt.Sprintf("`%s%s %s`\n%s", ctx.Conf.Prefix, command+sctx.Invoked, scalled.usage(), scalled.detailed())
			if len(scalled.subcommands()) != 0 {
				desc += "\n\nSubcommands:"
				desc += fmt.Sprintf(" `%shelp %s [subcommand]` for more info!", ctx.Conf.Prefix, command+sctx.Invoked)
				ckeys := make([]string, 0, len(scalled.subcommands()))
				for key := range scalled.subcommands() {
					ckeys = append(ckeys, key)
				}
				sc := sort.StringSlice(ckeys)
				sc.Sort()
				for _, k := range sc {
					desc += fmt.Sprintf("\n`%s%s %s` - %s", ctx.Conf.Prefix, command, k, scalled.subcommands()[k].description())
				}
			}
		} else {
			desc = "No command called `" + command + "` found!"
		}
	} else {
		desc = "Commands:"
		desc += fmt.Sprintf(" `%shelp [command]` for more info!", ctx.Conf.Prefix)
		keys := make([]string, 0, len(ch.Categories))
		for k := range ch.Categories {
			keys = append(keys, k)
		}
		s := sort.StringSlice(keys)
		s.Sort()
		for _, k := range s {
			var fdesc string
			field := &discordgo.MessageEmbedField{Name: k + ":"}
			ckeys := make([]string, 0, len(ch.Categories[k]))
			for key := range ch.Categories[k] {
				ckeys = append(ckeys, key)
			}
			sc := sort.StringSlice(ckeys)
			sc.Sort()
			for _, n := range ckeys {
				fdesc += fmt.Sprintf("\n`%s%s` - %s", ctx.Conf.Prefix, n, ch.Categories[k][n].description())
			}
			field.Value = fdesc[1:]
			embed.Fields = append(embed.Fields, field)
		}
	}
	embed.Author = &discordgo.MessageEmbedAuthor{Name: ctx.Mess.Author.Username, IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", ctx.Mess.Author.ID, ctx.Mess.Author.Avatar)}
	embed.Description = desc
	embed.Description += "\n\n"
	embed.Description += versionMarkdown()
	ctx.SendEm(embed)
}

func versionMarkdown() (versionMarkdown string) {
	return "[FloSelfbot](https://github.com/Moonlington/FloSelfbot) " + version
}
