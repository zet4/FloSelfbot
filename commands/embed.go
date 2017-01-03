package commands

import (
	"fmt"
	"strings"
)

type Embed struct{}

func (e *Embed) Message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) != 0 {
		text := strings.Join(ctx.Args, " ")
		em.Description = fmt.Sprintf("%s", text)
		ctx.SendEm(em)
	} else {
		em.Description = fmt.Sprintf("***%s*** *was silent...*", ctx.Mess.Author.Username)
		ctx.SendEmNoDelete(em)
	}
}

func (e *Embed) Description() string             { return "Embeds stuff" }
func (e *Embed) Usage() string                   { return "<message>" }
func (e *Embed) Detailed() string                { return "Embeds stuff." }
func (e *Embed) Subcommands() map[string]Command { return make(map[string]Command) }
