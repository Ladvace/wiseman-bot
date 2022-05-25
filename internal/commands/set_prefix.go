package commands

import (
	"context"
	"fmt"
	"log"
	"wiseman/internal/db"
	"wiseman/internal/discord"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "setPrefix",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})

	discord.Commands["setprefix"] = Setprefix
}

func Setprefix(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {

	ctx := context.TODO()
	if len(args) == 0 {
		log.Println("Expected arguments")
		return nil
	}
	prefix := args[0]
	collection := db.SERVERS_DB
	server := db.GetServerByID(m.GuildID)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"serverid": m.GuildID},
		bson.D{
			primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "guildprefix", Value: prefix}}},
		},
	)

	server.ServerPrefix = prefix
	db.UpsertServerByID(m.GuildID, server)

	if err == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Prefix set to %#v", prefix))
	}

	return nil
}
