package commands

import (
	"fmt"
	"log"
	"wiseman/internal/db"
	"wiseman/internal/entities"
	"wiseman/internal/errors"
	"wiseman/internal/services"
	"wiseman/internal/shared"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "deleterank",
		Category:    "Administrator Commands",
		Description: "deleterank sets the range of levels a role can be assigned to a user",
		Usage:       "deleterank <rank_id>",
	})

	services.Commands["deleterank"] = DeleteRank
}

func DeleteRank(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {

	//check if the user has the required role
	if !services.IsUserAdmin(m.Author.ID, m.ChannelID) {
		return errors.CreateUnauthorizedUserError(m.Author.ID)
	}

	if len(args) != 1 {
		if len(args) < 1 {
			s.ChannelMessageSend(m.ChannelID, "Not enough arguments")
			return errors.CreateInvalidArgumentError(args[0])
		} else if len(args) > 1 {
			s.ChannelMessageSend(m.ChannelID, "Too many arguments")
			return errors.CreateInvalidArgumentError(args[0])
		}
	}

	rank_id := args[0]

	customRole := entities.CustomRanks{
		Id: rank_id,
	}

	log.Printf("role %s deleted:", rank_id)

	db.DeleteRoleServer(m.GuildID, rank_id)
	services.UpdateUsersRoles(m.GuildID, shared.DELETE_OP, customRole)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role %s deleted", rank_id))

	return nil
}
