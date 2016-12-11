package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

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
	Token  string
	Prefix string
}

func createConfig() *Config {
	var tempprefix string
	var temptoken string

	fmt.Println("No config file found, so let's make one!")
	fmt.Println("\nTo find your User Token. In browser or desktop Discord, type Ctrl-Shift-I. Go to the Console section, and type localStorage.token. Your user token will appear. Do not share this token with anyone! This token provides complete access to your Discord account, so never share it!")
	fmt.Println("\nInput your User Token here: ")
	fmt.Scanln(&temptoken)
	fmt.Println("Input your desired prefix here: ")
	fmt.Scanln(&tempprefix)

	buf := new(bytes.Buffer)
	tempconfig := &Config{temptoken, tempprefix}

	logwarning(toml.NewEncoder(buf).Encode(tempconfig))

	f, err := os.Create("config.toml")
	logwarning(err)
	defer f.Close()

	_, err = f.Write(buf.Bytes())

	return tempconfig
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
		conf = createConfig()
	}

	dg, err := discordgo.New(conf.Token)

	logwarning(err)

	dg.AddHandler(messageCreate)
	commandhandler = &CommandHandler{make(map[string]Command)}

	commandhandler.AddCommand("ping", &Ping{})
	commandhandler.AddCommand("setgame", &SetGame{})
	commandhandler.AddCommand("me", &Me{})
	commandhandler.AddCommand("eval", &Eval{})
	commandhandler.AddCommand("clean", &Clean{})
	commandhandler.AddCommand("quote", &Quote{})
	commandhandler.AddCommand("afk", &Afk{})

	err = dg.Open()

	logwarning(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	fmt.Println("Type", conf.Prefix+"help", "to see all commands!")

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

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

	if strings.HasPrefix(m.Content, conf.Prefix) {
		// Setting values for the commands
		args := strings.Split(m.Content[len(conf.Prefix):len(m.Content)], " ")
		invoked := args[0]
		args = args[1:]
		channel, _ := s.State.Channel(m.ChannelID)
		guild, _ := s.State.Guild(channel.GuildID)

		ctx := &Context{invoked, args, channel, guild, m, s}

		commandhandler.HandleCommands(ctx)
	}
}
