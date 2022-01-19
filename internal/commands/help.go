package commands

import (
	"wiseman/internal/discord"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
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

func Help(s *discordgo.Session, m *discordgo.MessageCreate, mongo *mongo.Client) error {
	for _, v := range Helpers {
		s.ChannelMessageSend(m.ChannelID, v.Name)
	}

	return nil
}
