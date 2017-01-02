package commands

import (
	"fmt"
	"time"
)

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
