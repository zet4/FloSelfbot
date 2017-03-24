package commands

import (
	"io/ioutil"
	"net/http"
)

// Dog struct handles Dog Command
type Dog struct{}

func (m *Dog) message(ctx *Context) {
	resp, err := http.Get("http://random.dog/woof")
	if err != nil {
		em := createEmbed(ctx)
		em.Description = "Sorry, dog not present."
		ctx.SendEm(em)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	filename := string(body)

	if filename == "" {
		em := createEmbed(ctx)
		em.Description = "Sorry, dog missing."
		ctx.SendEm(em)
	}

	channelID := ctx.Channel.ID
	ctx.Sess.ChannelMessageSend(channelID, "http://random.dog/"+filename)
}

func (m *Dog) description() string             { return "Dog" }
func (m *Dog) usage() string                   { return "" }
func (m *Dog) detailed() string                { return "Dog." }
func (m *Dog) subcommands() map[string]Command { return make(map[string]Command) }
