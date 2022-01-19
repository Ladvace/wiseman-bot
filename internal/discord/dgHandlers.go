package discord

import (
	"fmt"
	"strings"
	"wiseman/internal/servers"

	"github.com/bwmarrin/discordgo"
)

type CommandFunc func(*discordgo.Session, *discordgo.MessageCreate) error

var Commands map[string]CommandFunc

func init() {
	Commands = make(map[string]CommandFunc, 200)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(servers.Get(m.GuildID).GuildPrefix, m.Content[0:1], Commands)

	// Check if prefix for this server is correct
	if servers.Get(m.GuildID).GuildPrefix != m.Content[0:1] {
		return
	}

	msg := strings.Split(m.Content[1:], " ")[0]

	// Check if command exists
	if _, ok := Commands[msg]; !ok {
		return
	}

	err := Commands[msg](s, m)
	if err != nil {
		fmt.Println(err)
		return
	}
}
