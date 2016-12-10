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

var conf *Config
var chandler *CommandHandler

func check(e error) {
	if e != nil {
		log.Fatal(e)
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

	check(toml.NewEncoder(buf).Encode(tempconfig))

	f, err := os.Create("config.toml")
	check(err)
	defer f.Close()

	_, err = f.Write(buf.Bytes())

	return tempconfig
}

func main() {
	_, err := toml.DecodeFile("config.toml", &conf)
	if os.IsNotExist(err) {
		conf = createConfig()
	}

	dg, err := discordgo.New(conf.Token)

	dg.AddHandler(messageCreate)
	chandler = &CommandHandler{make(map[string]Command)}

	chandler.AddCommand("ping", &Ping{})
	chandler.AddCommand("setgame", &SetGame{})
	chandler.AddCommand("me", &Me{})

	err = dg.Open()

	check(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	fmt.Println("Your prefix is", conf.Prefix)

	<-make(chan struct{})
	return
}

func getUserColor(s *discordgo.Session, guild *discordgo.Guild, userID string) int {

	u, err := s.GuildMember(guild.ID, userID)
	if err != nil {
		return 0
	}
	var highestrole *discordgo.Role
	highestrole = &discordgo.Role{Position: 0}

	for _, role := range guild.Roles {
		for _, roleid := range u.Roles {
			if role.ID == roleid {
				if role.Color != 0 {
					if highestrole.Position < role.Position {
						highestrole = role
					}
				}
			}
		}
	}
	return highestrole.Color
}

type Context struct {
	Invoked string
	Args    []string
	Channel *discordgo.Channel
	Guild   *discordgo.Guild
	Message *discordgo.MessageCreate
	Sess    *discordgo.Session
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

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

		chandler.HandleCommands(ctx)
	}
}
