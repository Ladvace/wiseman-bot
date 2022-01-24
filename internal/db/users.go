package db

import (
	"context"
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserType struct {
	ComplexID              string `bson:"complexid"`
	ServerID               string `bson:"serverid"`
	UserID                 string `bson:"userid"`
	MessagesCount          uint   `bson:"messagescount"`
	CurrentLevelExperience uint   `bson:"currentlevelexperience"`
	CurrentLevel           uint   `bson:"currentlevel"`
	LastTimeOnline         uint64 `bson:"lastranktime"`
	Bot                    bool   `bson:"bot"`
	Verified               bool   `bson:"verified"`
}

type UsersType map[string]UserType

func (u UserType) GetNextLevelMinExperience() uint {
	user := users[u.ComplexID]
	fLevel := float64(user.CurrentLevel + 1)

	return uint(50 * (math.Pow(fLevel, 3) - 6*math.Pow(fLevel, 2) + 17*fLevel - 12) / 3)
}

func (u UserType) IncreaseExperience(v uint) uint {
	// Get original object using ComplexID to avoid injecting other mutated data
	user := users[u.ComplexID]

	for {
		if user.CurrentLevelExperience+v < user.GetNextLevelMinExperience() {
			user.CurrentLevelExperience += v
			UpsertUserByID(u.ComplexID, user)
			break
		}

		v -= user.GetNextLevelMinExperience() - user.CurrentLevelExperience
		user.CurrentLevelExperience = 0
		user.CurrentLevel += 1
		UpsertUserByID(u.ComplexID, user)
	}

	r, _ := USERS_DB.UpdateOne(context.TODO(),
		bson.M{
			"complexid": u.ComplexID,
		},
		bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{
					Key: "currentlevel", Value: user.CurrentLevel,
				},
				primitive.E{
					Key: "currentlevelexperience", Value: user.CurrentLevelExperience,
				},
			}},
		},
	)

	fmt.Printf("%+v\n", r)

	return u.CurrentLevelExperience
}

var users UsersType

var USERS_DB *mongo.Collection

func init() {
	users = make(map[string]UserType, 50000)
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
				var user UserType
				err := res.Decode(&user)
				if err != nil {
					return 0, err
				}

				UpsertUserByID(memberID, user)
				continue
			}

			fmt.Println("User not found in DB", memberID, member.User.Username+"#"+member.User.Discriminator)
			nu += 1

			user := UserType{
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

func GetUserByID(userID, guildID string) UserType {
	return users[userID+"|"+guildID]
}

func UpsertUserByID(userID string, user UserType) {
	users[userID] = user
}
