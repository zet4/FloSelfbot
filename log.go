package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
)

func LogMessage(s *discordgo.Session, timestamp time.Time, user *discordgo.User, mID, cID, code, message string) {

	timestampo := timestamp.Format("2006-01-02 15:04:05") + " UTC"

	var namestr string

	channel, _ := s.State.Channel(cID)
	member, err := s.State.Member(channel.GuildID, user.ID)
	if err != nil {
		namestr = user.Username + "#" + user.Discriminator
	} else {
		if member.Nick != "" {
			namestr = member.Nick + " " + "(" + user.Username + "#" + user.Discriminator + ")"
		} else {
			namestr = user.Username + "#" + user.Discriminator
		}
	}

	InsertInLogFile(s, cID, fmt.Sprintf("%s %s %s %s ## %s ## %s\n", mID, timestampo, user.ID, code, namestr, message))
}

func InsertInLogFile(s *discordgo.Session, ChannelID, str string) {
	var channel *discordgo.Channel
	var guildname string
	var channelname string

	channel, _ = s.State.Channel(ChannelID)
	guild, err := s.State.Guild(channel.GuildID)

	if err != nil {
		guildname = "Direct Message"
		channelname = channel.Recipient.Username
	} else {
		guildname = guild.ID
		channelname = channel.Name
	}
	f, err := GetLogFile(guildname, guild.Name, channelname)
	logerror(err)
	defer f.Close()
	f.WriteString(str)
}

func GetLogFile(g, gn, c string) (*os.File, error) {
	os.MkdirAll(filepath.Join("logs", g), os.ModePerm)

	_, err := os.Stat(filepath.Join("logs", g, "servername.txt"))
	if os.IsNotExist(err) {
		f, err := os.Create(filepath.Join("logs", g, "servername.txt"))
		logerror(err)
		defer f.Close()
		f.WriteString(gn)
	}

	path := filepath.Join("logs", g, c+".txt")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return os.Create(path)
	}

	return os.OpenFile(path, os.O_APPEND, os.ModePerm)
}
