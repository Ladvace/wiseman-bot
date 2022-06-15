package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"wiseman/internal/db"
	"wiseman/internal/services"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "setmultiplier",
		Category:    "Points",
		Description: "setmultiplier sets the ",
		Usage:       "This is a usage",
	})

	services.Commands["setmultiplier"] = Setmultiplier
}

func Setmultiplier(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	ctx := context.TODO()
	if len(args) == 0 || len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Not enough arguments passed!")
		return nil
	}
	collection := db.SERVERS_DB
	multiplierType := strings.ToLower(args[0])
	multiplier := args[1]

	parsedMultiplied, err := strconv.ParseFloat(multiplier, 8)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Not right paramater passed")
		return nil
	}

	switch multiplierType {
	case "time":
		{
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"serverid": m.GuildID},
				bson.D{
					primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "timeexpmultiplier", Value: parsedMultiplied}}},
				},
			)
			if err != nil {
				fmt.Println("err", err)
				s.ChannelMessageSend(m.ChannelID, "Error while setting the time multiplier")
				return nil
			}

			server := db.GetServerByID(m.GuildID)
			server.TimeExpMultiplier = parsedMultiplied
			db.UpsertServerByID(m.GuildID, server)

			s.ChannelMessageSend(m.ChannelID, "Multiplier updated successfully!")
		}
	case "msg":
		{
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"serverid": m.GuildID},
				bson.D{
					primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "msgexpmultiplier", Value: parsedMultiplied}}},
				},
			)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error while setting the message multiplier")
				return nil
			}

			server := db.GetServerByID(m.GuildID)
			server.MsgExpMultiplier = parsedMultiplied
			db.UpsertServerByID(m.GuildID, server)

			s.ChannelMessageSend(m.ChannelID, "Multiplier updated successfully!")
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Not right type passed")
		return nil
	}

	return nil
}
