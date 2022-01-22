package commands

import (
	"wiseman/internal/discord"

	"github.com/bwmarrin/discordgo"
)

type Helper struct {
	Name        string
	Category    string
	Description string
	Usage       string
}

var Helpers []Helper

func init() {
	discord.Commands["help"] = Help
}

func Help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	for _, v := range Helpers {
		s.ChannelMessageSend(m.ChannelID, v.Name)
	}

	return nil
}
