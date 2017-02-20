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

var version = "compiled manually"

var currentgame string

// A Config struct contains variables used by Commands
type Config struct {
	Token              string
	Prefix             string
	LogMode            bool
	LogModeMinBuffer   int
	LogModeMaxBuffer   int
	LogModeCompression bool
	EmbedColor         string
	AFKPlay            bool
	MultiGameStrings   []string
	MultiGameMinutes   int
	MultigameToggled   bool
	AutoDeleteSeconds  int
	SketchyMode        bool
}

// A Context struct holds variables for Messages
type Context struct {
	Conf    *Config
	Invoked string
	Argstr  string
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

// EditConfigFile edits the config.toml file with new info
// conf: Config struct to edit it to
func EditConfigFile(conf *Config) {
	f, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	logwarning(err)
	defer f.Close()
	logwarning(toml.NewEncoder(f).Encode(conf))
}

func waitandDelete(ctx *Context, m *discordgo.Message) {
	time.Sleep(time.Second * time.Duration(ctx.Conf.AutoDeleteSeconds))
	ctx.Sess.ChannelMessageDelete(m.ChannelID, m.ID)
}

// SendEm is a helper function to easily send embeds
// em: embed to send
func (ctx *Context) SendEm(em *discordgo.MessageEmbed) (*discordgo.Message, error) {
	m, err := ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
	if ctx.Conf.AutoDeleteSeconds != 0 {
		go waitandDelete(ctx, m)
	}
	return m, err
}

// QuickSendEm is a helper function to easily send strings as an embed
// s: string to send
func (ctx *Context) QuickSendEm(s string) (*discordgo.Message, error) {
	em := createEmbed(ctx)
	em.Description = s
	return ctx.SendEm(em)
}

// SendEmNoDelete is a helper function to easily send embeds but doesn't use Autodelete
// em: embed to send
func (ctx *Context) SendEmNoDelete(em *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func createEmbed(ctx *Context) *discordgo.MessageEmbed {
	if ctx.Conf.EmbedColor == "#000000" || ctx.Conf.EmbedColor == "" {
		if ctx.Conf.EmbedColor == "" {
			ctx.Conf.EmbedColor = "#000000"
			EditConfigFile(ctx.Conf)
		}
		color := ctx.Sess.State.UserColor(ctx.Mess.Author.ID, ctx.Mess.ChannelID)
		return &discordgo.MessageEmbed{Color: color}
	}
	color := ctx.Conf.EmbedColor
	if strings.HasPrefix(color, "#") {
		color = "0x" + ctx.Conf.EmbedColor[1:]
	}
	d, _ := strconv.ParseInt(color, 0, 64)
	return &discordgo.MessageEmbed{Color: int(d)}
}
