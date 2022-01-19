package commands

import (
	"wiseman/internal/discord"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "setPrefix",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})

	discord.Commands["setPrefix"] = Setprefix
}

func Setprefix(s *discordgo.Session, m *discordgo.MessageCreate, mongo *mongo.Client) error {

	return nil
}
