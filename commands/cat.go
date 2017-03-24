package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type catpicStruct struct {
	Link string `json:"file"`
}

var catpic = catpicStruct{}

// Cat struct handles Cat Command
type Cat struct{}

func (m *Cat) message(ctx *Context) {
	resp, err := http.Get("http://random.cat/meow")
	if err != nil {
		em := createEmbed(ctx)
		em.Description = "Sorry, cat not work."
		ctx.SendEm(em)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &catpic)
	if err != nil || catpic.Link == "" {
		em := createEmbed(ctx)
		em.Description = "Sorry, cat broke."
		ctx.SendEm(em)
	}

	channelID := ctx.Channel.ID
	ctx.Sess.ChannelMessageSend(channelID, catpic.Link)
}

func (m *Cat) description() string             { return "Cat" }
func (m *Cat) usage() string                   { return "" }
func (m *Cat) detailed() string                { return "Cat." }
func (m *Cat) subcommands() map[string]Command { return make(map[string]Command) }
