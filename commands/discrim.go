package commands

import "fmt"

// Discrim struct handles Discrim Command
type Discrim struct{}

func (d *Discrim) message(ctx *Context) {
	users, err := ctx.GetAllUsers()
	var discrim string
	if err != nil {
		return
	}
	if len(ctx.Args) != 0 {
		discrim = ctx.Args[0]
	} else {
		discrim = ctx.Mess.Author.Discriminator
	}
	var matching string
	var i int
	for _, u := range users {
		if u.Discriminator == discrim {
			matching += u.Username + ", "
			i++
		}
	}
	if len(matching) != 0 {
		ctx.QuickSendEm(fmt.Sprintf("Found %d users with discrim `%s`:\n%s", i, discrim, matching[:len(matching)-2]))
	} else {
		ctx.QuickSendEm(fmt.Sprintf("No users found with discrim `%s`!", discrim))
	}
}

func (d *Discrim) description() string             { return "Finds users with your/specified discrim" }
func (d *Discrim) usage() string                   { return "[discrim]" }
func (d *Discrim) detailed() string                { return "Finds users with your/specified discrim" }
func (d *Discrim) subcommands() map[string]Command { return make(map[string]Command) }
