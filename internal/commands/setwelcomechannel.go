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
		Name:        "SetWelcomechannel",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})

	discord.Commands["setwelcomechannel"] = Setwelcomechannel
}

func Setwelcomechannel(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {

	ctx := context.TODO()
	if len(args) == 0 {
		return nil
	}
	channelId := args[0]
	collection := db.SERVERS_DB
	server := db.GetServerByID(m.GuildID)

	if channelId == "null" {
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"serverid": m.GuildID},
			bson.D{
				primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "welcomechannel", Value: ""}}},
			},
		)

		server.WelcomeChannel = ""
		db.UpsertServerByID(m.GuildID, server)

		if err == nil {
			s.ChannelMessageSend(m.ChannelID, "Welcome Channel has been reset!")
		}
	}

	channel, err := s.Channel(channelId)
	if err != nil {
		return nil
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"serverid": m.GuildID},
		bson.D{
			primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "welcomechannel", Value: channelId}}},
		},
	)

	server.WelcomeChannel = channelId
	db.UpsertServerByID(m.GuildID, server)

	if err == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Welcome Channel set to %#v", channel.Name))
	}

	return nil
}
