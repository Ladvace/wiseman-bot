package commands

import (
	"wiseman/internal/db"
	"wiseman/internal/errors"
	"wiseman/internal/services"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "resetranks",
		Category:    "Points",
		Description: "resetranks resets all ranks to their default values",
		Usage:       "resetranks",
	})

	services.Commands["resetranks"] = ResetRanks
}

func ResetRanks(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	//check if the user has the required role
	if !services.IsUserAdmin(m.Author.ID, m.ChannelID) {
		return errors.CreateUnauthorizedUserError(m.Author.ID)
	}

	err := db.ResetRanks()
	if err != nil {
		return err
	}

	s.ChannelMessageSend(m.ChannelID, "Ranks reset")

	return nil
}
