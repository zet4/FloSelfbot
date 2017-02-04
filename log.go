package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var logbuffers map[string]map[string]bytes.Buffer = make(map[string]map[string]bytes.Buffer)
var logmintime time.Time = time.Now().Add(time.Duration(conf.LogModeMinBuffer) * time.Minute)
var logmaxtime time.Time = time.Now().Add(time.Duration(conf.LogModeMaxBuffer) * time.Minute)

func SendToBuffer(s *discordgo.Session, ChannelID, str string) {
	var channel *discordgo.Channel
	var guildname string
	var channelname string

	channel, _ = s.State.Channel(ChannelID)
	guild, err := s.State.Guild(channel.GuildID)

	now := time.Now()

	logbuffer := logbuffers[guild.ID][ChannelID]

	if now.Before(logmintime) {
		logbuffer.WriteString(str)
	} else {
		var f *os.File

		if err != nil {
			guildname = "Direct Message"
			channelname = channel.Recipient.Username
			f, err = GetLogFile(guildname, "Direct Message", channelname)
		} else {
			guildname = guild.ID
			channelname = channel.Name
			f, err = GetLogFile(guildname, guild.Name, channelname)
		}
		logbuffer.WriteString(str)
		logbuffer.WriteTo(f)
		logbuffer.Reset()
	}
}

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

	SendToBuffer(s, cID, strings.Replace(fmt.Sprintf("%s %s %s %s ## %s ## %s", mID, timestampo, user.ID, code, namestr, message), "\n", "\t", -1)+"\n")
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
	re := regexp.MustCompile(`[\\/:\?!\*"<>\|]`)
	c = re.ReplaceAllString(c, "")
	path := filepath.Join("logs", g, c+".txt")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return os.Create(path)
	}

	return os.OpenFile(path, os.O_APPEND, os.ModePerm)
}
