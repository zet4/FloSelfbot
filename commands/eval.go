package commands

import (
	"fmt"
	"strings"

	"github.com/robertkrimen/otto"
)

// Eval struct handles Eval Command
type Eval struct{}

func (e *Eval) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		vm := otto.New()
		vm.Set("ctx", ctx)
		vm.Set("ctx.Conf", ctx.Conf)
		toEval := strings.Join(ctx.Args, " ")
		executed, err := vm.Run(toEval)
		em := createEmbed(ctx)
		if err != nil {
			em.Description = fmt.Sprintf("Input: `%s`\n\nError: `%s`", toEval, err.Error())
			ctx.SendEm(em)
			return
		}
		em.Description = fmt.Sprintf("Input: `%s`\n\nOutput: ```js\n%s\n```", toEval, executed.String())
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "You didn't eval anything"
		ctx.SendEm(em)
	}
}

func (e *Eval) description() string { return "Evaluates using Otto (Advanced stuff, don't bother)" }
func (e *Eval) usage() string       { return "<toEval>" }
func (e *Eval) detailed() string {
	return "Evaluates using Otto (Advanced stuff, don't bother)\nIf you do want to bother, the `ctx` variable is something you can use."
}
func (e *Eval) subcommands() map[string]Command { return make(map[string]Command) }
