package discord

import (
	"strings"
	"wiseman/internal/commands"
	"wiseman/internal/servers"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if prefix for this server is correct
	if servers.Get(m.GuildID).GuildPrefix != m.Content[0:1] {
		return
	}

	msg := strings.Split(m.Content[1:], " ")[0]

	switch msg {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "help":
		commands.Help(s, m)
	}
}
