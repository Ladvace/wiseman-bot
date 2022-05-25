package commands

import (
	"log"
	"strconv"
	"wiseman/internal/db"
	"wiseman/internal/discord"
	"wiseman/internal/entities"
	"wiseman/internal/errors"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "setrank",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})

	discord.Commands["setrank"] = SetRank
}

func SetRank(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {

	//check if the user has the required role
	if !discord.IsUserAdmin(m.Author.ID, m.ChannelID) {
		return errors.CreateUnauthorizedUserError(m.Author.ID)
	}

	// ctx := context.TODO()
	// rank_id min_xp max_xp
	if len(args) != 3 {
		log.Println("Expected arguments")
		return nil
	}

	rank_id := args[0]
	min_xp, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.CreateInvalidArgumentError(args[1])
	}

	max_xp, err := strconv.Atoi(args[2])
	if err != nil {
		return errors.CreateInvalidArgumentError(args[2])
	}

	customRole := &entities.RoleType{
		Id:       rank_id,
		MinLevel: uint(min_xp),
		MaxLevel: uint(max_xp),
	}

	log.Println("new role created:", customRole)

	db.UpdateRoleServer(m.GuildID, *customRole)
	// now the role exists, what we need to do is add it to the server
	// maybe with s.GuildRoleCreate(m.GuildId) and reallocate all
	// users who match this role to the new role

	return nil
}
