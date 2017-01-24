package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Embed struct{}

func (e *Embed) ParseEmbed(em *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	output := em.Description
	iterations := 0
	var lastoutput string
	var regexed []string
	for lastoutput != output && iterations < 200 {
		lastoutput = output
		i1 := strings.Index(output, "}")
		var i2 int
		if i1 == -1 {
			i2 = -1
		} else {
			i2 = strings.LastIndex(output[:i1], "{")
		}
		if i1 != -1 && i2 != -1 {
			var toEval = output[i2+1 : i1]
			if strings.HasPrefix(toEval, "author:") {
				emex := &discordgo.MessageEmbedAuthor{}
				regexed = regexp.MustCompile(`\[url:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.URL = regexed[1]
				}
				regexed = regexp.MustCompile(`\[iconurl:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.IconURL = regexed[1]
				}
				regexed = regexp.MustCompile(`\[name:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Name = regexed[1]
				}
				output = output[:i2] + output[i1+1:]
				em.Author = emex

			} else if strings.HasPrefix(toEval, "field:") {
				emex := &discordgo.MessageEmbedField{}
				inline := regexp.MustCompile(`\[inline\]`).Match([]byte(toEval))
				emex.Inline = inline
				regexed = regexp.MustCompile(`\[text:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Value = regexed[1]
				}
				regexed = regexp.MustCompile(`\[name:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Name = regexed[1]
				}
				output = output[:i2] + output[i1+1:]
				em.Fields = append(em.Fields, emex)

			} else if strings.HasPrefix(toEval, "footer:") {
				emex := &discordgo.MessageEmbedFooter{}
				regexed = regexp.MustCompile(`\[iconurl:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.IconURL = regexed[1]
				}
				regexed = regexp.MustCompile(`\[text:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Text = regexed[1]
				}
				output = output[:i2] + output[i1+1:]
				em.Footer = emex

			} else if strings.HasPrefix(toEval, "image:") {
				emex := &discordgo.MessageEmbedImage{}
				emex.URL = toEval[6:]
				output = output[:i2] + output[i1+1:]
				em.Image = emex

			} else if strings.HasPrefix(toEval, "provider:") {
				emex := &discordgo.MessageEmbedProvider{}
				regexed = regexp.MustCompile(`\[url:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.URL = regexed[1]
				}
				regexed = regexp.MustCompile(`\[name:(.*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Name = regexed[1]
				}
				output = output[:i2] + output[i1+1:]
				em.Provider = emex

			} else if strings.HasPrefix(toEval, "thumbnail:") {
				emex := &discordgo.MessageEmbedThumbnail{}
				emex.URL = toEval[10:]
				output = output[:i2] + output[i1+1:]
				em.Thumbnail = emex
			}
		}
	}
	em.Description = strings.TrimSpace(lastoutput)
	return em
}
func (e *Embed) Message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) != 0 {
		text := strings.Join(ctx.Args, " ")
		em.Description = text
		ctx.SendEmNoDelete(e.ParseEmbed(em))
	} else {
		em.Description = fmt.Sprintf("***%s*** *was silent...*", ctx.Mess.Author.Username)
		ctx.SendEmNoDelete(em)
	}
}

func (e *Embed) Description() string             { return "Embeds stuff" }
func (e *Embed) Usage() string                   { return "<message>" }
func (e *Embed) Detailed() string                { return "Embeds stuff." }
func (e *Embed) Subcommands() map[string]Command { return make(map[string]Command) }
