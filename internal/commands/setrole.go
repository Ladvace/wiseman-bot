package commands

import (
	"context"
	"fmt"
	"wiseman/internal/db"
	"wiseman/internal/discord"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "setRole",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})

	discord.Commands["setrole"] = Setrole
}

func Setrole(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	ctx := context.TODO()
	if len(args) == 0 {
		return nil
	}
	level := args[0]
	roleID := args[1]
	collection := db.SERVERS_DB
	server := db.GetServerByID(m.GuildID)

	roles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		return nil
	}

	var roleName string
	for _, role := range roles {
		if role.ID == roleID {
			roleName = role.Name
		}
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"serverid": m.GuildID},
		bson.D{
			primitive.E{Key: "$set", Value: bson.M{fmt.Sprintf("customranks.%#v", level): roleID}},
		},
	)

	str := string(level)
	server.CustomRanks[str] = roleID
	db.UpsertServerByID(m.GuildID, server)

	if err == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role %#v set at level %#v", roleName, level))
	}

	return nil
}
