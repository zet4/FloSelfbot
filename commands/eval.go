package commands

import (
	"fmt"
	"strings"

	"github.com/robertkrimen/otto"
)

type Eval struct{}

func (e *Eval) Message(ctx *Context) {
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

func (e *Eval) Description() string { return "Evaluates using Otto (Advanced stuff, don't bother)" }
func (e *Eval) Usage() string       { return "<toEval>" }
func (e *Eval) Detailed() string {
	return "Evaluates using Otto (Advanced stuff, don't bother)\nIf you do want to bother, the `ctx` variable is something you can use."
}
func (e *Eval) Subcommands() map[string]Command { return make(map[string]Command) }
