package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func removeDuplicateUsers(list *[]*discordgo.User) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *list {
		if !found[x.ID] {
			found[x.ID] = true
			(*list)[j] = (*list)[i]
			j++
		}
	}
	*list = (*list)[:j]
}

// GetAllUsers is a helper function to return all members
func (ctx *Context) GetAllUsers() (ms []*discordgo.User, err error) {
	servers, err := ctx.Sess.UserGuilds(100, "", "")
	if err != nil {
		return
	}
	for _, server := range servers {
		var g *discordgo.Guild
		g, err = ctx.Sess.State.Guild(server.ID)
		if err != nil {
			continue
		}
		for _, m := range g.Members {
			ms = append(ms, m.User)
		}
	}
	removeDuplicateUsers(&ms)
	return
}

// GetUserByName is a helper function to find a User by string
// query: String to use when finding User
func (ctx *Context) GetUserByName(query string) (members []*discordgo.User, err error) {
	MentionRegex := regexp.MustCompile(`<@!?(\d+)>`)
	var id, discrim string
	if MentionRegex.MatchString(query) {
		id = MentionRegex.FindStringSubmatch(query)[1]
	} else if regexp.MustCompile(`^.*#\d{4}$`).MatchString(query) {
		discrim = query[len(query)-4:]
		query = strings.TrimSpace(query[:len(query)-5])
	}
	var exact, wrongcase, startswith, contains, all []*discordgo.User
	lowerQuery := strings.ToLower(query)
	all, err = ctx.GetAllUsers()
	if err != nil {
		return
	}
	for _, u := range all {
		if id != "" && u.ID == id {
			exact = append(exact, u)
			break
		}
		if discrim != "" && u.Discriminator != discrim {
			continue
		}
		if u.Username == query {
			exact = append(exact, u)
		} else if len(exact) == 0 && strings.ToLower(u.Username) == lowerQuery {
			wrongcase = append(wrongcase, u)
		} else if len(wrongcase) == 0 && strings.HasPrefix(strings.ToLower(u.Username), lowerQuery) {
			startswith = append(startswith, u)
		} else if len(startswith) == 0 && strings.Contains(strings.ToLower(u.Username), lowerQuery) {
			contains = append(contains, u)
		}
	}
	if len(exact) != 0 {
		members = exact
	} else if len(wrongcase) != 0 {
		members = wrongcase
	} else if len(startswith) != 0 {
		members = startswith
	} else {
		members = contains
	}
	return
}

// GuildGetUserByName is a helper function to find a User by string in a guild
// query: String to use when finding User
// GuildID: ID for guild
func (ctx *Context) GuildGetUserByName(query, GuildID string) (members []*discordgo.User, err error) {
	MentionRegex := regexp.MustCompile(`<@!?(\d+)>`)
	var id, discrim string
	if MentionRegex.MatchString(query) {
		id = MentionRegex.FindStringSubmatch(query)[1]
	} else if regexp.MustCompile(`^.*#\d{4}$`).MatchString(query) {
		discrim = query[len(query)-4:]
		query = strings.TrimSpace(query[:len(query)-5])
	}
	var exact, wrongcase, startswith, contains []*discordgo.User
	var all []*discordgo.Member
	lowerQuery := strings.ToLower(query)
	g, err := ctx.Sess.State.Guild(GuildID)
	if err != nil {
		return
	}
	all = g.Members
	if err != nil {
		return
	}
	for _, m := range all {
		u := m.User
		if id != "" && u.ID == id {
			exact = append(exact, u)
			break
		}
		if discrim != "" && u.Discriminator != discrim {
			continue
		}
		if u.Username == query {
			exact = append(exact, u)
		} else if len(exact) == 0 && strings.ToLower(u.Username) == lowerQuery {
			wrongcase = append(wrongcase, u)
		} else if len(wrongcase) == 0 && strings.HasPrefix(strings.ToLower(u.Username), lowerQuery) {
			startswith = append(startswith, u)
		} else if len(startswith) == 0 && strings.Contains(strings.ToLower(u.Username), lowerQuery) {
			contains = append(contains, u)
		}
	}
	if len(exact) != 0 {
		members = exact
	} else if len(wrongcase) != 0 {
		members = wrongcase
	} else if len(startswith) != 0 {
		members = startswith
	} else {
		members = contains
	}
	return
}

// ParseTooManyUsers is a helper function to create a message for finding too many users
// query: String used when finding User
// list: List of users found
func (ctx *Context) ParseTooManyUsers(query string, list []*discordgo.User) (*discordgo.Message, error) {
	out := fmt.Sprintf("Multiple users found for query **%s**:", query)
	for i := 0; i < 6; i++ {
		if i < len(list) {
			out += "\n - " + list[i].Username + " #" + list[i].Discriminator
		}
	}
	if len(list) > 6 {
		out += "\n**And " + strconv.Itoa(len(list)-6) + " more...**"
	}
	return ctx.QuickSendEm(out)
}
