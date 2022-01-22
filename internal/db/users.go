package db

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserType struct {
	ComplexID      string `bson:"complexid"`
	ServerID       string `bson:"serverid"`
	UserID         string `bson:"userid"`
	MessagesCount  uint   `bson:"messagescount"`
	Rank           int    `bson:"rank"`
	Time           uint   `bson:"time"`
	Experience     uint   `bson:"experience"`
	LastTimeOnline uint64 `bson:"lastranktime"`
	Bot            bool   `bson:"bot"`
	Verified       bool   `bson:"verified"`
}

type UsersType map[string]UserType

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
				ComplexID:      memberID,
				UserID:         member.User.ID,
				ServerID:       v.ServerID,
				Verified:       member.User.Verified,
				Bot:            member.User.Bot,
				MessagesCount:  0,
				Rank:           0,
				Time:           0,
				Experience:     0,
				LastTimeOnline: 0,
			}

			USERS_DB.InsertOne(context.TODO(), user)
			UpsertUserByID(memberID, user)
		}
	}

	return nu, nil
}

func GetUserByID(userID string) UserType {
	return users[userID]
}

func UpsertUserByID(userID string, user UserType) {
	users[userID] = user
}
