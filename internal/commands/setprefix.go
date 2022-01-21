package commands

import (
	"context"
	"fmt"
	"wiseman/internal/db"
	"wiseman/internal/discord"
	"wiseman/internal/shared"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func Setprefix(s *discordgo.Session, m *discordgo.MessageCreate, client *mongo.Client, args []string) error {

	ctx := context.TODO()
	if len(args) == 0 {
		return nil
	}
	prefix := args[0]
	collection := client.Database(shared.DB_NAME).Collection(shared.SERVERS_INFIX)
	server := db.GetServerById(m.GuildID)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"serverid": m.GuildID},
		bson.D{
			primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "guildprefix", Value: prefix}}},
		},
	)

	server.ServerPrefix = prefix
	db.UpsertServerById(m.GuildID, server)

	if err == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Prefix set to %#v", prefix))
	}

	return nil
}
