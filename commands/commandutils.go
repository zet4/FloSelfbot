package commands

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

const version string = "3.4"

var currentgame string = ""

type Config struct {
	Token             string
	Prefix            string
	LogMode           bool
	LogModeMinBuffer  int
	LogModeMaxBuffer  int
	EmbedColor        string
	AFKPlay           bool
	MultiGameStrings  []string
	MultiGameMinutes  int
	MultigameToggled  bool
	AutoDeleteSeconds int
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

func WaitandDelete(ctx *Context, m *discordgo.Message) {
	time.Sleep(time.Second * time.Duration(ctx.Conf.AutoDeleteSeconds))
	ctx.Sess.ChannelMessageDelete(m.ChannelID, m.ID)
}

func (ctx *Context) SendEm(em *discordgo.MessageEmbed) (*discordgo.Message, error) {
	m, err := ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
	if ctx.Conf.AutoDeleteSeconds != 0 {
		go WaitandDelete(ctx, m)
	}
	return m, err
}

func (ctx *Context) QuickSendEm(s string) (*discordgo.Message, error) {
	em := createEmbed(ctx)
	em.Description = s
	return ctx.SendEm(em)
}

func (ctx *Context) SendEmNoDelete(em *discordgo.MessageEmbed) (*discordgo.Message, error) {
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
