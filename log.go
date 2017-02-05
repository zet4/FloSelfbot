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

var logbuffers map[string]map[string]*bytes.Buffer = make(map[string]map[string]*bytes.Buffer)
var logmintime time.Time
var logmaxtime time.Time

func SendToBuffer(s *discordgo.Session, ChannelID, str string) {
	var channel *discordgo.Channel
	var gn string
	var cn string

	channel, _ = s.State.Channel(ChannelID)
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		gn = "Direct Message"
		cn = channel.Recipient.Username + "#" + channel.Recipient.Discriminator
	} else {
		gn = guild.ID
		cn = channel.Name
	}

	now := time.Now()

	logbuffer, ok := logbuffers[gn][cn]

	if !ok {
		if _, kok := logbuffers[gn]; !kok {
			logbuffers[gn] = make(map[string]*bytes.Buffer)
		}
		buf := new(bytes.Buffer)
		logbuffers[gn][cn] = buf
		logbuffer = buf
	}

	if now.Before(logmintime) {
		logbuffer.WriteString(str)
	} else {
		logbuffer.WriteString(str)
		if logbuffer.Len() >= logbuffer.Cap() {
			f, err := GetLogFile(s, gn, cn)
			logerror(err)
			defer f.Close()
			n, err := logbuffer.WriteTo(f)
			logerror(err)
			fmt.Println("Written", n, "bytes to", filepath.Join(filepath.Join("logs", gn, cn+".txt")))
		}
	}
}

func BufferLoop(s *discordgo.Session) {
	if conf.LogModeMaxBuffer < 1 {
		conf.LogModeMaxBuffer = 5
		editConfigfile(conf)
	}
	if conf.LogModeMinBuffer < 1 {
		conf.LogModeMinBuffer = 10
		editConfigfile(conf)
	}
	logmaxtime = time.Now().Add(time.Duration(conf.LogModeMaxBuffer) * time.Minute)
	logmintime = time.Now().Add(time.Duration(conf.LogModeMinBuffer) * time.Minute)
	for {
		if time.Now().After(logmaxtime) {
			tn := 0
			for k, v := range logbuffers {
				for c, buf := range v {
					f, err := GetLogFile(s, k, c)
					logerror(err)
					defer f.Close()
					n, err := buf.WriteTo(f)
					logerror(err)
					tn += int(n)
				}
			}
			fmt.Println("Flushed buffers, written ", tn, "bytes")
			logmintime = time.Now().Add(time.Duration(conf.LogModeMinBuffer) * time.Minute)
			logmaxtime = time.Now().Add(time.Duration(conf.LogModeMaxBuffer) * time.Minute)
		}
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

func GetLogFile(s *discordgo.Session, g, c string) (*os.File, error) {
	os.MkdirAll(filepath.Join("logs", g), os.ModePerm)
	if g != "Direct Message" {
		_, err := os.Stat(filepath.Join("logs", g, "servername.txt"))
		if os.IsNotExist(err) {
			guild, err := s.State.Guild(g)
			f, err := os.Create(filepath.Join("logs", g, "servername.txt"))
			logerror(err)
			defer f.Close()

			f.WriteString(guild.Name)
		}
	}
	re := regexp.MustCompile(`[\\/:\?!\*"<>\|]`)
	c = re.ReplaceAllString(c, "")
	path := filepath.Join("logs", g, c+".txt")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.Create(path)
	}

	return os.OpenFile(path, os.O_APPEND, os.ModePerm)
}
