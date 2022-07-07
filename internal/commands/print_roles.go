package commands

import (
	"fmt"
	"wiseman/internal/db"
	"wiseman/internal/services"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "printranks",
		Category:    "Administrator Commands",
		Description: "printranks shows active ranks on the server",
		Usage:       "printranks",
	})

	services.Commands["printranks"] = printranks
}

func printranks(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {

	c := db.GetCustomRanksByGuildId(m.GuildID)
	if len(c) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No available roles found")
	}
	for _, r := range c {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Available roles are:\nRole Id: %v - Role Min Level: %d - Role Max Level: %d", r.Id, r.MinLevel, r.MaxLevel))
	}

	return nil
}
