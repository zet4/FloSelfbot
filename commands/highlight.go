package commands

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var highlightRegexCache []*regexp.Regexp
var highlightRegexCacheMutex sync.Mutex

const highlightFormat = "**[$0](.)**"

func highlightRebuildCache(conf *Config) {
	highlightRegexCacheMutex.Lock()
	highlightRegexCache = make([]*regexp.Regexp, len(conf.HighlightStrings))
	for idx, keyword := range conf.HighlightStrings {
		highlightRegexCache[idx] = regexp.MustCompile(keyword)
	}
	highlightRegexCacheMutex.Unlock()
}

// HighLightFunc handles Highlight
func HighLightFunc(s *discordgo.Session, conf *Config, m *discordgo.Message) {
	if len(highlightRegexCache) != len(conf.HighlightStrings) {
		highlightRebuildCache(conf)
	}
	var newMsg string
	matches := make([]string, 0)
	for _, highlight := range highlightRegexCache {
		if localMatches := highlight.FindAllStringIndex(m.Content, -1); localMatches != nil {
			resultsArr := make([]string, len(localMatches))
			if len(newMsg) == 0 {
				newMsg = highlight.ReplaceAllString(m.Content, highlightFormat)
			} else {
				newMsg = highlight.ReplaceAllString(newMsg, highlightFormat)
			}
			for idx, match := range localMatches {
				start, end := match[0], match[1]
				resultsArr[idx] = m.Content[start:end] + " (" + strconv.Itoa(start) + ":" + strconv.Itoa(end) + ")"
			}
			matches = append(matches, "`"+highlight.String()+"`: "+strings.Join(resultsArr, ", "))
		}
	}
	if len(matches) == 0 {
		return
	}

	ch, err := s.State.Channel(m.ChannelID)
	var authorIcon, guildIcon string
	var guild *discordgo.Guild
	if err == nil && !ch.IsPrivate {
		guild, _ = s.State.Guild(ch.GuildID)
		if len(guild.Icon) > 0 {
			guildIcon = discordgo.EndpointGuildIcon(guild.ID, guild.Icon)
		}
	}
	authorIcon = discordgo.EndpointUserAvatar(m.Author.ID, m.Author.Avatar)
	emauthor := &discordgo.MessageEmbedAuthor{Name: m.Author.Username, IconURL: authorIcon}
	timestamp, err := m.Timestamp.Parse()
	logerror(err)
	timestampo := timestamp.Local().Format(time.ANSIC)
	sentfrom := "Sent "
	if ch.IsPrivate {
		sentfrom = sentfrom + "in DM with " + ch.Recipient.Username + "#" + ch.Recipient.Discriminator
	} else {
		sentfrom = sentfrom + "from #" + ch.Name + " in " + guild.Name
	}

	emfooter := &discordgo.MessageEmbedFooter{Text: sentfrom + " | " + timestampo, IconURL: guildIcon}
	emcolor := s.State.UserColor(m.Author.ID, m.ChannelID)
	em := &discordgo.MessageEmbed{Author: emauthor, Footer: emfooter, Description: newMsg[0:int(math.Min(2000, float64(len(newMsg))))], Color: emcolor}
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Matches", Value: strings.Join(matches, "\n"), Inline: false})
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Owner", Value: "<@" + m.Author.ID + ">", Inline: true})
	if !ch.IsPrivate {
		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Channel", Value: "<#" + ch.ID + ">", Inline: true})
	}
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Message ID", Value: m.ID, Inline: true})
	split := strings.Split(conf.HighlightWebhook, "/")
	err = s.WebhookExecute(split[5], split[6], false, &discordgo.WebhookParams{Embeds: []*discordgo.MessageEmbed{em}})
	logerror(err)
}

type addHighLightString struct{}

func (a *addHighLightString) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		keyword := ctx.Argstr
		if strings.ToLower(keyword) == "all" {
			ctx.QuickSendEm("**All** can not be added because of problems with the remove command.")
			return
		} else if _, err := regexp.Compile(keyword); err != nil {
			ctx.QuickSendEm(err.Error())
			return
		}
		em := createEmbed(ctx)
		ctx.Conf.HighlightStrings = append(ctx.Conf.HighlightStrings, keyword)
		highlightRebuildCache(ctx.Conf)
		EditConfigFile(ctx.Conf)
		em.Description = fmt.Sprintf("Added **`%s`** to Highlight list", keyword)
		ctx.SendEm(em)
	} else {
		em := createEmbed(ctx)
		em.Description = "No keyword specified!"
		ctx.SendEm(em)
	}
}

func (a *addHighLightString) description() string { return "Adds a string to the Highlight." }
func (a *addHighLightString) usage() string       { return "<keyword>" }
func (a *addHighLightString) detailed() string {
	return "Adds a string to the Highlight list."
}
func (a *addHighLightString) subcommands() map[string]Command { return make(map[string]Command) }

type removeHighLightString struct{}

func (r *removeHighLightString) message(ctx *Context) {
	if len(ctx.Args) != 0 {
		var pos int
		em := createEmbed(ctx)
		keyword := ctx.Argstr
		if strings.ToLower(keyword) == "all" {
			ctx.Conf.HighlightStrings = []string{}
			highlightRebuildCache(ctx.Conf)
			EditConfigFile(ctx.Conf)
			em.Description = "Removed everything from the list."
			ctx.SendEm(em)
			return
		}
		for i, v := range ctx.Conf.HighlightStrings {
			if keyword == v {
				pos = i
				break
			} else {
				pos = -1
			}
		}
		if pos == -1 {
			em.Description = "Game `" + keyword + "` not found in Highlight! (Check for caps and stuff.)"
			ctx.SendEm(em)
		} else {
			ctx.Conf.HighlightStrings = append(ctx.Conf.HighlightStrings[:pos], ctx.Conf.HighlightStrings[pos+1:]...)
			highlightRebuildCache(ctx.Conf)
			EditConfigFile(ctx.Conf)
			em.Description = fmt.Sprintf("Removed **%s** from Highlight", keyword)
			ctx.SendEm(em)
		}
	} else {
		em := createEmbed(ctx)
		em.Description = "No keyword name specified!"
		ctx.SendEm(em)
	}
}

func (r *removeHighLightString) description() string { return "Removes a string from Highlight" }
func (r *removeHighLightString) usage() string       { return "<keyword|all>" }
func (r *removeHighLightString) detailed() string {
	return "Removes a string from Highlight. Type All as keyword to remove everything from the list"
}
func (r *removeHighLightString) subcommands() map[string]Command { return make(map[string]Command) }

type highLightList struct{}

func (l *highLightList) message(ctx *Context) {
	var desc string
	em := createEmbed(ctx)
	if len(ctx.Conf.HighlightStrings) == 0 {
		desc = "There are currently no keywords in Highlight!"
	} else {
		desc = "Current strings in Highlight:"
		for _, v := range ctx.Conf.HighlightStrings {
			desc += fmt.Sprintf("\n`%s`", v)
		}
	}
	em.Description = desc
	ctx.SendEm(em)
}

func (l *highLightList) description() string { return "Returns your Highlight strings" }
func (l *highLightList) usage() string       { return "" }
func (l *highLightList) detailed() string {
	return "Returns the list of all your strings in Highlight"
}
func (l *highLightList) subcommands() map[string]Command { return make(map[string]Command) }

type highLightToggle struct{}

func (mgt *highLightToggle) message(ctx *Context) {
	if ctx.Conf.HighlightWebhook == "" {
		if resp, err := ctx.Sess.WebhookCreate(ctx.Channel.ID, "Highlight", ""); err != nil {
			logwarning(err)
		} else {
			ctx.Conf.HighlightWebhook = discordgo.EndpointWebhookToken(resp.ID, resp.Token)
		}
	} else {
		// dirty hack
		split := strings.Split(ctx.Conf.HighlightWebhook, "/")
		if _, err := ctx.Sess.WebhookDeleteWithToken(split[5], split[6]); err != nil {
			logwarning(err)
		}
		ctx.Conf.HighlightWebhook = ""

	}
	highlightRebuildCache(ctx.Conf)
	EditConfigFile(ctx.Conf)

	em := createEmbed(ctx)
	em.Description = fmt.Sprintf("Toggled HighLight to **%s**", strconv.FormatBool(ctx.Conf.HighlightWebhook != ""))
	ctx.SendEm(em)
}

func (mgt *highLightToggle) description() string { return "Toggles Highlight" }
func (mgt *highLightToggle) usage() string       { return "" }
func (mgt *highLightToggle) detailed() string {
	return "Toggles if highlight is on or off."
}
func (mgt *highLightToggle) subcommands() map[string]Command { return make(map[string]Command) }

// HighLight struct handles HighLight Command
type HighLight struct{}

func (mg *HighLight) message(ctx *Context) {
	em := createEmbed(ctx)
	if len(ctx.Args) == 0 {
		em.Description = "Command `highlight` requires a subcommand!"
	} else {
		em.Description = fmt.Sprintf(`Subcommand "%s" doesn't exist!`, ctx.Args[0])
	}
	ctx.SendEm(em)
}

func (mg *HighLight) description() string { return `Commands for Highlight file` }
func (mg *HighLight) usage() string       { return "" }
func (mg *HighLight) detailed() string {
	return "Commands related to highlight, the webhook and regexp based notification system."
}
func (mg *HighLight) subcommands() map[string]Command {
	return map[string]Command{"add": &addHighLightString{}, "remove": &removeHighLightString{}, "list": &highLightList{}, "toggle": &highLightToggle{}}
}
