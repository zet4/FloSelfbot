package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var regionalindicators = []string{"🇦", "🇧", "🇨", "🇩", "🇪", "🇫", "🇬", "🇭", "🇮", "🇯", "🇰", "🇱", "🇲", "🇳", "🇴", "🇵", "🇶", "🇷", "🇸", "🇹", "🇺", "🇻", "🇼", "🇽", "🇾", "🇿"}

// Poll struct handles Poll Command
type Poll struct{}

func (p *Poll) message(ctx *Context) {
	if len(ctx.Args) > 0 {
		var toreact []string
		em := createEmbed(ctx)
		em.Author = &discordgo.MessageEmbedAuthor{IconURL: discordgo.EndpointUserAvatar(ctx.Mess.Author.ID, ctx.Mess.Author.Avatar), Name: ctx.Mess.Author.Username}
		split := strings.Split(strings.Trim(ctx.Argstr, ":"), ":")
		var desc string
		if len(split) > 1 {
			if len(split) > 21 {
				ctx.QuickSendEm("Too many choices! (Max is 20)")
				return
			}
			desc += split[0] + "\n"
			for i, choice := range split[1:] {
				desc += fmt.Sprintf("\n%s %s", regionalindicators[i], strings.TrimSpace(choice))
				toreact = append(toreact, regionalindicators[i])
			}
		} else {
			desc += ctx.Argstr
			toreact = append(toreact, "👍", "👎")
		}
		em.Description = desc
		msg, _ := ctx.SendEmNoDelete(em)
		for _, r := range toreact {
			ctx.Sess.MessageReactionAdd(msg.ChannelID, msg.ID, r)
		}
	} else {
		ctx.QuickSendEm("No question specified!")
	}
}

func (p *Poll) description() string { return "Creates a poll" }
func (p *Poll) usage() string       { return "<question>:[choice]:[choice]..." }
func (p *Poll) detailed() string {
	return "Creates a poll.\nIf you don't specify choices, it defaults to a `YES/NO` question, else it will be a `A B C D...` question.\nChoices must be split with `:`"
}
func (p *Poll) subcommands() map[string]Command { return make(map[string]Command) }
