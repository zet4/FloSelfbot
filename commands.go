package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Ping struct{}

func (p *Ping) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
	color := getUserColor(ctx.Sess, ctx.Guild, ctx.Message.Author.ID)
	start := time.Now()
	msg, _ := ctx.Sess.ChannelMessageSendEmbed(ctx.Message.ChannelID, &discordgo.MessageEmbed{Description: "Pong!", Color: color})
	elapsed := time.Since(start)
	ctx.Sess.ChannelMessageEditEmbed(ctx.Message.ChannelID, msg.ID, &discordgo.MessageEmbed{Description: fmt.Sprintf("Pong! `%s`", elapsed), Color: color})
}

func (p *Ping) Description() string { return "Measures latency" }

type SetGame struct{}

func (sg *SetGame) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
	color := getUserColor(ctx.Sess, ctx.Guild, ctx.Message.Author.ID)
	game := strings.Join(ctx.Args, " ")
	ctx.Sess.UpdateStatus(0, game)
	ctx.Sess.ChannelMessageSendEmbed(ctx.Message.ChannelID, &discordgo.MessageEmbed{Description: fmt.Sprintf("Changed game to: **%s**", game), Color: color})

}

func (p *SetGame) Description() string { return "Sets your game to anything you like" }

type Me struct{}

func (m *Me) Message(ctx *Context) {
	ctx.Sess.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
	color := getUserColor(ctx.Sess, ctx.Guild, ctx.Message.Author.ID)
	text := strings.Join(ctx.Args, " ")
	ctx.Sess.ChannelMessageSendEmbed(ctx.Message.ChannelID, &discordgo.MessageEmbed{Description: fmt.Sprintf("***%s*** *%s*", ctx.Message.Author.Username, text), Color: color})
}

func (p *Me) Description() string { return "Says stuff" }
