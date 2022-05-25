package db

import (
	"context"
	"fmt"
	"strings"
	"wiseman/internal/entities"

	"github.com/bwmarrin/discordgo"
	"github.com/r3labs/diff/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var users entities.UsersType
var USERS_DB *mongo.Collection

func init() {
	users = make(map[string]entities.UserType, 50000)
}

func UpdateExpById(userID, guildID string, exp int) error {

	userStruct := users[userID+"|"+guildID]

	userStruct.CurrentLevelExperience += uint(exp)
	UpdateUser(userID, guildID, userStruct)
}

func GetUserByID(userID, guildID string) entities.UserType {
	return users[userID+"|"+guildID]
}

func UpdateUser(userId, guildId string, u entities.UserType) {
	users[userId+"|"+guildId] = u
}

func HydrateUsers(d *discordgo.Session) (int, error) {
	var nu int
	for _, v := range servers {
		var members []*discordgo.Member
		var lastID string
		for {
			newMembers, err := d.GuildMembers(v.ServerID, lastID, 1000)
			if err != nil {
				return 0, err
			}

			if len(newMembers) == 0 {
				break
			}

			lastID = newMembers[len(newMembers)-1].User.ID
			members = append(members, newMembers...)
		}

		// TODO: Use InsertMany to optimize this
		for _, member := range members {
			memberID := member.User.ID + "|" + v.ServerID
			// Check if server is already in DB
			res := USERS_DB.FindOne(context.TODO(), bson.M{"complexid": memberID})
			if res.Err() != mongo.ErrNoDocuments {
				var user entities.UserType
				err := res.Decode(&user)
				if err != nil {
					return 0, err
				}

				UpsertUserByID(memberID, user)
				continue
			}

			fmt.Println("User not found in DB", memberID, member.User.Username+"#"+member.User.Discriminator)
			nu += 1

			user := entities.UserType{
				ComplexID:              memberID,
				UserID:                 member.User.ID,
				ServerID:               v.ServerID,
				Verified:               member.User.Verified,
				Bot:                    member.User.Bot,
				MessagesCount:          0,
				CurrentLevelExperience: 0,
				CurrentLevel:           1,
				LastTimeOnline:         0,
			}

			USERS_DB.InsertOne(context.TODO(), user)
			UpsertUserByID(memberID, user)
		}
	}

	return nu, nil
}

func UpsertUserByID(userID string, user entities.UserType) {
	if Hydrated {
		d, err := diff.NewDiffer(diff.TagName("bson"))
		if err != nil {
			panic(err)
		}
		changelog, err := d.Diff(users[userID], user)
		if err != nil {
			panic(err)
		}
		if len(changelog) == 0 {
			return
		}

		changes := bson.D{}
		for _, v := range changelog {
			changes = append(changes, primitive.E{
				Key: "$set",
				Value: bson.D{
					primitive.E{
						Key:   strings.Join(v.Path, "."),
						Value: v.To,
					},
				},
			})
		}
		_, err = USERS_DB.UpdateOne(
			context.TODO(),
			bson.M{"complexid": user.ComplexID},
			changes,
		)

		if err != nil {
			panic(err)
		}
	}

	users[userID] = user
}
