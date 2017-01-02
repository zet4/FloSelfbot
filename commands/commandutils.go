package commands

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

const version string = "3.1"

type Config struct {
	Token            string
	Prefix           string
	LogMode          bool
	EmbedColor       string
	MultiGameStrings []string
	MultiGameMinutes int
	MultigameToggled bool
}

type Context struct {
	Conf    *Config
	Invoked string
	Args    []string
	Channel *discordgo.Channel
	Guild   *discordgo.Guild
	Mess    *discordgo.MessageCreate
	Sess    *discordgo.Session
}

func logwarning(e error) {
	if e != nil {
		log.Println(e)
	}
}

func logerror(e error) {
	if e != nil {
		log.Println(e)
	}
}

func editConfigfile(conf *Config) {
	f, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	logwarning(err)
	defer f.Close()
	logwarning(toml.NewEncoder(f).Encode(conf))
}

func (ctx *Context) SendEm(em *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func createEmbed(ctx *Context) *discordgo.MessageEmbed {
	if ctx.Conf.EmbedColor == "#000000" || ctx.Conf.EmbedColor == "" {
		if ctx.Conf.EmbedColor == "" {
			ctx.Conf.EmbedColor = "#000000"
			editConfigfile(ctx.Conf)
		}
		color := ctx.Sess.State.UserColor(ctx.Mess.Author.ID, ctx.Mess.ChannelID)
		return &discordgo.MessageEmbed{Color: color}
	} else {
		color := ctx.Conf.EmbedColor
		if strings.HasPrefix(color, "#") {
			color = "0x" + ctx.Conf.EmbedColor[1:]
		}
		d, _ := strconv.ParseInt(color, 0, 64)
		return &discordgo.MessageEmbed{Color: int(d)}
	}
}
