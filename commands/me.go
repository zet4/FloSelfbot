package commands

import (
	"fmt"
)

// Me struct handles Me Command
type Me struct{}

func (m *Me) message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) != 0 {
		text := ctx.Argstr
		em.Description = fmt.Sprintf("***%s*** *%s*", ctx.Mess.Author.Username, text)
		ctx.SendEmNoDelete(em)
	} else {
		em.Description = fmt.Sprintf("***%s*** *was silent...*", ctx.Mess.Author.Username)
		ctx.SendEmNoDelete(em)
	}
}

func (m *Me) description() string             { return "Says stuff" }
func (m *Me) usage() string                   { return "<message>" }
func (m *Me) detailed() string                { return "Says stuff." }
func (m *Me) subcommands() map[string]Command { return make(map[string]Command) }
