package db

import (
	"context"
	"fmt"
	"wiseman/internal/shared"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserType struct {
	ComplexID     string `bson:"complexid"`
	serverID      string `bson:"serverid"`
	UserID        string `bson:"userid"`
	Username      string `bson:"username"`
	Discriminator string `bson:"discriminator"`
	MessagesCount uint   `bson:"messagescount"`
	Rank          int    `bson:"rank"`
	Time          uint   `bson:"time"`
	Exp           uint   `bson:"exp"`
	LastRankTime  uint64 `bson:"lastranktime"`
	Bot           bool   `bson:"bot"`
	Verified      bool   `bson:"verified"`
}

type UsersType map[string]UserType

func init() {
	users = make(map[string]UserType, 50000)
}

func HydrateUsers(d *discordgo.Session, m *mongo.Client) (int, error) {
	var nu int
	for _, v := range servers {
		var members []*discordgo.Member
		var lastId string
		for {
			newMembers, err := d.GuildMembers(v.ServerID, lastId, 1000)
			if err != nil {
				return 0, err
			}

			if len(newMembers) == 0 {
				break
			}

			lastId = newMembers[len(newMembers)-1].User.ID
			members = append(members, newMembers...)
		}

		// TODO: Use InsertMany to optimize this
		for _, member := range members {
			memberId := member.User.ID + "#" + member.User.Discriminator + "|" + v.ServerID
			// Check if server is already in DB
			res := m.Database(shared.DB_NAME).Collection(shared.USERS_INFIX).FindOne(context.TODO(), bson.M{"complexid": memberId})
			if res.Err() != mongo.ErrNoDocuments {
				var user UserType
				err := res.Decode(&user)
				if err != nil {
					return 0, err
				}
				UpsertUserById(memberId, user)
				continue
			}

			fmt.Println("User not found in DB", memberId, member.User.Username+"#"+member.User.Discriminator)
			nu += 1

			user := UserType{
				ComplexID:     memberId,
				UserID:        member.User.ID,
				serverID:      v.ServerID,
				Username:      member.User.Username,
				Discriminator: member.User.Discriminator,
				Verified:      member.User.Verified,
				Bot:           member.User.Bot,
				MessagesCount: 0,
				Rank:          0,
				Time:          0,
				Exp:           0,
				LastRankTime:  0,
			}

			m.Database(shared.DB_NAME).Collection(shared.USERS_INFIX).InsertOne(context.TODO(), user)
			UpsertUserById(memberId, user)
		}
	}

	return nu, nil
}

func GetUserById(userId string) UserType {
	return users[userId]
}

func UpsertUserById(userId string, user UserType) {
	users[userId] = user
}
