package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Embed struct{}

func (e *Embed) ParseEmbed(em *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	output := strings.Replace(strings.Replace(em.Description, `\\]`, "\u0014", -1), `\\[`, "\u0013", -1)
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
				regexed = regexp.MustCompile(`\[url:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.URL = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
				}
				regexed = regexp.MustCompile(`\[iconurl:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.IconURL = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
				}
				regexed = regexp.MustCompile(`\[name:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Name = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
				}
				output = output[:i2] + output[i1+1:]
				em.Author = emex

			} else if strings.HasPrefix(toEval, "field:") {
				emex := &discordgo.MessageEmbedField{}
				inline := regexp.MustCompile(`\[inline\]`).Match([]byte(toEval))
				emex.Inline = inline
				regexed = regexp.MustCompile(`\[text:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Value = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
				}
				regexed = regexp.MustCompile(`\[name:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Name = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
				}
				output = output[:i2] + output[i1+1:]
				em.Fields = append(em.Fields, emex)

			} else if strings.HasPrefix(toEval, "footer:") {
				emex := &discordgo.MessageEmbedFooter{}
				regexed = regexp.MustCompile(`\[iconurl:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.IconURL = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
				}
				regexed = regexp.MustCompile(`\[text:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Text = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
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
				regexed = regexp.MustCompile(`\[url:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.URL = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
				}
				regexed = regexp.MustCompile(`\[name:([\s\S]*?)\]`).FindStringSubmatch(toEval)
				if regexed != nil {
					emex.Name = strings.Replace(strings.Replace(regexed[1], "\u0014", `]`, -1), "\u0013", `[`, -1)
				}
				output = output[:i2] + output[i1+1:]
				em.Provider = emex

			} else if strings.HasPrefix(toEval, "thumbnail:") {
				emex := &discordgo.MessageEmbedThumbnail{}
				emex.URL = toEval[10:]
				output = output[:i2] + output[i1+1:]
				em.Thumbnail = emex
			} else if strings.HasPrefix(toEval, "color:") {
				color := toEval[6:]
				colorint, err := strconv.Atoi(color)
				logerror(err)
				output = output[:i2] + output[i1+1:]
				em.Color = colorint
			}
		}
	}
	em.Description = strings.Replace(strings.Replace(strings.TrimSpace(lastoutput), "\u0014", `]`, -1), "\u0013", `[`, -1)
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
