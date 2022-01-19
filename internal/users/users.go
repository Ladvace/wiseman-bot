package users

import (
	"context"
	"fmt"
	"wiseman/internal/servers"
	"wiseman/internal/shared"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserType struct {
	UserId        string `bson:"userid"`
	MessagesCount uint   `bson:"messagescount"`
	Rank          int    `bson:"rank"`
	Time          uint   `bson:"time"`
	Exp           uint   `bson:"exp"`
	GuildId       string `bson:"guildid"`
	LastRankTime  uint64 `bson:"lastranktime"`
	Bot           bool   `bson:"bot"`
	Verified      bool   `bson:"verified"`
}

type UsersType map[string]UserType

var users UsersType

func init() {
	users = make(map[string]UserType, 50000)
}

func GetAll() *UsersType {
	return &users
}

func Get(id string) UserType {
	return users[id]
}

func Upsert(id string, u UserType) {
	users[id] = u
}

func Hydrate(d *discordgo.Session, m *mongo.Client) error {

	for k, _ := range *servers.GetAll() {
		var members []*discordgo.Member
		var lastId string
		for {
			newMembers, err := d.GuildMembers(k, lastId, 1000)
			if err != nil {
				return err
			}

			if len(newMembers) == 0 {
				break
			}

			lastId = newMembers[len(newMembers)-1].User.ID
			members = append(members, newMembers...)
		}

		// TODO: Use InsertMany to optimize this
		for _, member := range members {
			// Check if server is already in DB
			res := m.Database(shared.DB_NAME).Collection(shared.USERS_INFIX).FindOne(context.TODO(), bson.M{"userid": member.User.ID})
			if res.Err() != mongo.ErrNoDocuments {
				var user UserType
				err := res.Decode(&user)
				if err != nil {
					return err
				}
				Upsert(member.User.ID, user)
				continue
			}

			fmt.Println("User not found in DB", member.User.ID, member.User.Username+"#"+member.User.Discriminator)
			user := UserType{
				UserId:        member.User.ID,
				Verified:      member.User.Verified,
				Bot:           member.User.Bot,
				MessagesCount: 0,
				Rank:          0,
				Time:          0,
				Exp:           0,
				GuildId:       k,
				LastRankTime:  0,
			}

			m.Database(shared.DB_NAME).Collection(shared.USERS_INFIX).InsertOne(context.TODO(), user)
			Upsert(member.User.ID, user)
		}
	}

	return nil
}
