package commands

import (
	"context"
	"fmt"
	"strings"
	"wiseman/internal/db"
	"wiseman/internal/discord"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "setWelcomeMsg",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})

	discord.Commands["setwelcomemsg"] = SetwelcomeMsg
}

func SetwelcomeMsg(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	ctx := context.TODO()
	if len(args) == 0 {
		return nil
	}
	welcomeMessage := strings.Join(args, " ")
	collection := db.SERVERS_DB
	server := db.GetServerByID(m.GuildID)

	if welcomeMessage == "null" {
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"serverid": m.GuildID},
			bson.D{
				primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "welcomemessage", Value: ""}}},
			},
		)

		server.WelcomeMessage = ""
		db.UpsertServerByID(m.GuildID, server)
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, "Welcome message has been reset!")
		}
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"serverid": m.GuildID},
		bson.D{
			primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "welcomemessage", Value: welcomeMessage}}},
		},
	)

	server.WelcomeMessage = welcomeMessage
	db.UpsertServerByID(m.GuildID, server)

	if err == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Welcome message set to %#v", welcomeMessage))
	}

	return nil
}
