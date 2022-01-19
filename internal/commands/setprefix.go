package commands

import (
	"wiseman/internal/discord"

	"github.com/bwmarrin/discordgo"
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

func Setprefix(s *discordgo.Session, m *discordgo.MessageCreate) error {

	return nil
}
