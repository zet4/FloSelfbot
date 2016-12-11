package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robertkrimen/otto"
)

func createEmbed(ctx *Context) *discordgo.MessageEmbed {
	color := ctx.Sess.State.UserColor(ctx.Mess.Author.ID, ctx.Mess.ChannelID)
	return &discordgo.MessageEmbed{Color: color}
}

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

type SetGame struct{}

func (sg *SetGame) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	em := createEmbed(ctx)
	game := strings.Join(ctx.Args, " ")
	em.Description = game
	ctx.Sess.UpdateStatus(0, game)
	ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)

}

func (p *SetGame) Description() string { return "Sets your game to anything you like" }

type Me struct{}

func (m *Me) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.ID)
	em := createEmbed(ctx)
	text := strings.Join(ctx.Args, " ")
	em.Description = fmt.Sprintf("***%s*** *%s*", ctx.Mess.Author.Username, text)
	ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func (m *Me) Description() string { return "Says stuff" }

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
