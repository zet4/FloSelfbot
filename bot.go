package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

const version string = "3.0"

var (
	conf           *Config
	commandhandler *CommandHandler
	Error          *log.Logger
	Warning        *log.Logger
	AFKMode        bool
	AFKMessages    []*discordgo.MessageCreate
	AFKstring      string
)

func logwarning(e error) {
	if e != nil {
		Warning.Println(e)
	}
}

func logerror(e error) {
	if e != nil {
		Error.Println(e)
	}
}

type Config struct {
	Token            string
	Prefix           string
	LogMode          bool
	EmbedColor       string
	MultiGameStrings []string
	MultiGameMinutes int
	MultigameToggled bool
}

func editConfigfile(conf *Config) {
	f, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	logwarning(err)
	defer f.Close()
	logwarning(toml.NewEncoder(f).Encode(conf))
}

func createConfig() *Config {
	var (
		tempprefix string
		temptoken  string
	)

	fmt.Println("\nTo find your User Token. In browser or desktop Discord, type Ctrl-Shift-I. Go to the Console section, and type localStorage.token. Your user token will appear. Do not share this token with anyone! This token provides complete access to your Discord account, so never share it!")
	fmt.Print("\nInput your User Token here: ")
	fmt.Scanln(&temptoken)
	fmt.Print("\nInput your desired prefix here: ")
	fmt.Scanln(&tempprefix)

	tempconfig := &Config{temptoken, tempprefix, false, "#000000", []string{}, 5, false}
	editConfigfile(tempconfig)

	return tempconfig
}

func MultiGameFunc(s *discordgo.Session) {
	for {
		if len(conf.MultiGameStrings) != 0 && conf.MultigameToggled {
			a := conf.MultiGameStrings
			newstring, a := a[0], a[1:]
			conf.MultiGameStrings = append(a, newstring)
			err := s.UpdateStatus(0, newstring)
			logerror(err)
		}
		if conf.MultiGameMinutes < 1 {
			conf.MultiGameMinutes = 1
			editConfigfile(conf)
		}
		time.Sleep(time.Minute * time.Duration(conf.MultiGameMinutes))
	}
}

func main() {

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	AFKMode = false

	_, err := toml.DecodeFile("config.toml", &conf)
	if os.IsNotExist(err) {
		fmt.Println("No config file found, so let's make one!")
		conf = createConfig()
	}

	dg, err := discordgo.New(conf.Token)

	logwarning(err)

	dg.AddHandler(messageCreate)
	commandhandler = &CommandHandler{make(map[string]Command)}

	commandhandler.AddCommand("ping", &Ping{})
	commandhandler.AddCommand("setgame", &SetGame{})
	commandhandler.AddCommand("me", &Me{})
	commandhandler.AddCommand("embed", &Embed{})
	commandhandler.AddCommand("eval", &Eval{})
	commandhandler.AddCommand("clean", &Clean{})
	commandhandler.AddCommand("quote", &Quote{})
	commandhandler.AddCommand("afk", &Afk{})
	commandhandler.AddCommand("config", &Configcommand{})
	commandhandler.AddCommand("multigame", &MultiGame{})
	// commandhandler.AddCommand("emote", &Emote{})

	err = dg.Open()

	logwarning(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	fmt.Println("Type", conf.Prefix+"help", "to see all commands!")

	go MultiGameFunc(dg)

	<-make(chan struct{})
	return
}

type Context struct {
	Invoked string
	Args    []string
	Channel *discordgo.Channel
	Guild   *discordgo.Guild
	Mess    *discordgo.MessageCreate
	Sess    *discordgo.Session
}

func (ctx *Context) SendEm(em *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Sess.ChannelMessageSendEmbed(ctx.Mess.ChannelID, em)
}

func logmessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if conf.LogMode {
		var f *os.File
		var channel *discordgo.Channel
		var guildname string
		var channelname string

		channel, _ = s.State.Channel(m.ChannelID)
		guild, err := s.State.Guild(channel.GuildID)
		if err != nil {
			guildname = "Direct Message"
			channelname = channel.Recipient.Username
		} else {
			guildname = guild.Name
			channelname = channel.Name
		}

		f, err = os.OpenFile(fmt.Sprintf("./logs/%s/%s.txt", guildname, channelname), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if os.IsNotExist(err) {
			os.MkdirAll(fmt.Sprintf("./logs/%s", guildname), 0777)
			f, _ = os.OpenFile(fmt.Sprintf("./logs/%s/%s.txt", guildname, channelname), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		}
		defer f.Close()

		timestamp, err := m.Timestamp.Parse()
		logerror(err)

		timestampo := timestamp.Format(time.ANSIC)
		f.Write([]byte(fmt.Sprintf("%s %s#%s (%s): %s\r\n", timestampo, m.Author.Username, m.Author.Discriminator, m.Author.ID, m.ContentWithMentionsReplaced())))
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	logmessage(s, m)

	if AFKMode {
		for _, u := range m.Mentions {
			if u.ID == s.State.User.ID {
				AFKMessages = append(AFKMessages, m)
				emcolor := s.State.UserColor(s.State.User.ID, m.ChannelID)
				em := &discordgo.MessageEmbed{Color: emcolor, Title: fmt.Sprintf("**%s** Is AFK!", s.State.User.Username)}
				if AFKstring != "" {
					em.Description = AFKstring
					s.ChannelMessageSendEmbed(m.ChannelID, em)
				} else {
					s.ChannelMessageSendEmbed(m.ChannelID, em)
				}
			}
		}
	}

	// Ignore all messages created by other users
	if m.Author.ID != s.State.User.ID {
		return
	}

	// i1 := strings.Index(m.Content, ":")
	// if i1 != -1 {
	// 	i2 := strings.Index(m.Content[i1+1:], ":")
	// 	if i2 != -1 {
	// 		emote := m.Content[i1+1 : i2+1]
	// 		exp, ok := conf.Emotes[emote]
	// 		if ok {
	// 			content := m.Content[:i1] + exp + m.Content[i2:]
	// 			s.ChannelMessageSend(m.ChannelID, content)
	// 		}
	// 	}
	// }

	if strings.HasPrefix(m.Content, conf.Prefix) {
		// Setting values for the commands
		var ctx *Context
		args := strings.Split(m.Content[len(conf.Prefix):len(m.Content)], " ")
		invoked := args[0]
		args = args[1:]
		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			channel, err = s.State.PrivateChannel(m.ChannelID)
			ctx = &Context{invoked, args, channel, nil, m, s}
		} else {
			guild, _ := s.State.Guild(channel.GuildID)
			ctx = &Context{invoked, args, channel, guild, m, s}
		}

		commandhandler.HandleCommands(ctx)
	}
}
