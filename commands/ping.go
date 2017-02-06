package commands

import (
	"fmt"
	"time"
)

// Ping struct handles Ping Command
type Ping struct{}

func (p *Ping) message(ctx *Context) {
	em := createEmbed(ctx)
	em.Description = "Pong!"
	start := time.Now()
	msg, _ := ctx.SendEm(em)
	elapsed := time.Since(start)
	em.Description = fmt.Sprintf("Pong! `%s`", elapsed)
	ctx.Sess.ChannelMessageEditEmbed(ctx.Mess.ChannelID, msg.ID, em)
}

func (p *Ping) description() string             { return "Measures latency" }
func (p *Ping) usage() string                   { return "" }
func (p *Ping) detailed() string                { return "Measures latency" }
func (p *Ping) subcommands() map[string]Command { return make(map[string]Command) }
