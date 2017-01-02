package commands

import (
	"fmt"
	"strings"
)

type Me struct{}

func (m *Me) Message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) != 0 {
		text := strings.Join(ctx.Args, " ")
		em.Description = fmt.Sprintf("***%s*** *%s*", ctx.Mess.Author.Username, text)
		ctx.SendEm(em)
	} else {
		em.Description = fmt.Sprintf("***%s*** *was silent...*", ctx.Mess.Author.Username)
		ctx.SendEm(em)
	}
}

func (m *Me) Description() string             { return "Says stuff" }
func (m *Me) Usage() string                   { return "<message>" }
func (m *Me) Detailed() string                { return "Says stuff." }
func (m *Me) Subcommands() map[string]Command { return make(map[string]Command) }
