package commands

import (
	"context"
	"fmt"
	"sort"
	"time"
	"wiseman/internal/db"
	"wiseman/internal/discord"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "leaderboard",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})

	discord.Commands["leaderboard"] = Leaderboard
}

type LeaderboardPlace struct {
	Level uint
	Field *discordgo.MessageEmbedField
}

func Leaderboard(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {

	ctx := context.TODO()
	collection := db.USERS_DB

	findOptions := options.Find()
	// Sort by `rank` field descending
	findOptions.SetSort(bson.D{primitive.E{Key: "rank", Value: -1}})
	// Limit by 10 documents only
	findOptions.SetLimit(10)

	cursor, err := collection.Find(ctx, bson.D{primitive.E{Key: "serverid", Value: m.GuildID}}, findOptions)
	if err != nil {
		return err
	}

	var fields []LeaderboardPlace

	for cursor.Next(ctx) {
		var leaderboard db.UserType
		err := cursor.Decode(&leaderboard)
		if err != nil {
			return err
		}

		leaderboardUser, err := discord.RetrieveUser(leaderboard.UserID, m.GuildID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fields = append(fields, LeaderboardPlace{
			Level: leaderboard.CurrentLevel,
			Field: &discordgo.MessageEmbedField{
				Name:  string(leaderboardUser.Username),
				Value: fmt.Sprint("Level ", leaderboard.CurrentLevel),
			}},
		)
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Level > fields[j].Level
	})

	finalFields := make([]*discordgo.MessageEmbedField, 0)

	for _, v := range fields {
		finalFields = append(finalFields, v.Field)
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       9004799,
		Description: "top 10 active users.",
		Fields:      finalFields,
		Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:       "Leaderboard",
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
