package commands

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// FindEmoji struct handles FindEmoji Command
type FindEmoji struct{}

func (r *FindEmoji) message(ctx *Context) {
	regex, err := regexp.Compile(ctx.Argstr)

	if err != nil {
		ctx.QuickSendEm("Error compiling your regex!")
	}

	em := createEmbed(ctx)
	em.Description = "Emojis on servers matching `" + ctx.Argstr + "`."

	servers := make(map[string][]string)
mainLoop:
	for _, guild := range ctx.Sess.State.Guilds {
		for _, emoji := range guild.Emojis {
			if len(emoji.Roles) == 0 && len(regex.FindStringSubmatch(emoji.APIName())) > 0 {
				if len(servers) >= 25 {
					em.Description = em.Description + "\n\n\tOver 25 results, try to narrow down your query."
					break mainLoop
				}
				if out := strings.Join(servers[guild.Name], " "); len(out) >= 1000 {
					em.Description = em.Description + "\n\t" + guild.Name + " has too many matches."
					break
				} else {
					servers[guild.Name] = append(servers[guild.Name], "<:z:"+emoji.ID+">")
				}
			}
		}
	}

	if len(servers) == 0 {
		em.Description = em.Description + "\n\n\tFound nothing..."
	} else {
		for guild := range servers {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{guild + " (" + strconv.Itoa(len(servers[guild])) + " total)", strings.Join(servers[guild], " "), true})
		}
	}

	_, err = ctx.SendEm(em)
	if err != nil {
		ctx.QuickSendEm(err.Error())
	}
}

func (r *FindEmoji) description() string             { return "Find emojis matching given regex" }
func (r *FindEmoji) usage() string                   { return "<Regex>" }
func (r *FindEmoji) detailed() string                { return "Find emojis matching given regex" }
func (r *FindEmoji) subcommands() map[string]Command { return make(map[string]Command) }
