package db

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"wiseman/internal/entities"

	"github.com/bwmarrin/discordgo"
	"github.com/r3labs/diff/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users struct {
	cache  map[string]*entities.UserType
	writes int
	lock   sync.RWMutex
}

var users Users = Users{
	cache:  make(map[string]*entities.UserType, 50000),
	writes: 0,
}

var USERS_DB *mongo.Collection

func UpdateExpById(userID, guildID string, exp int) {
	userStruct := users[userID+"|"+guildID]
	userStruct.CurrentLevelExperience += uint(exp)
	UpdateUser(userID, guildID, userStruct)
}

func UpdateUser(userId, guildId string, u *entities.UserType) {
	users[userId+"|"+guildId] = u
}

func HydrateUsers(d *discordgo.Session) (int, error) {
	var nu int
	for _, v := range servers.cache {
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

				UpsertUserByID(memberID, &user)
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
			UpsertUserByID(memberID, &user)
		}
	}

	return nu, nil
}

func GetUserByID(userID, guildID string) *entities.UserType {
	users.lock.RLock()
	u := users.cache[userID+"|"+guildID]
	users.lock.RUnlock()

	return u
}

func ResetRanks() error {
	for _, v := range users.cache {
		v.CurrentLevel = 1
		v.CurrentLevelExperience = 0
	}
	return nil
}

func UpsertUserByID(userID string, user *entities.UserType) error {
	if Hydrated {
		d, err := diff.NewDiffer(diff.TagName("bson"))
		if err != nil {
			return err
		}
		changelog, err := d.Diff(user, users.cache[userID])
		if err != nil {
			return err
		}
		if len(changelog) == 0 {
			return err
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
			return err
		}
	}

	if !Hydrated {
		users.lock.Lock()
		users.cache[userID] = user
		users.lock.Unlock()
	}

	return nil
}

func UpdateUser(userID string, user *entities.UserType) {
	users.lock.Lock()
	users.cache[userID] = user
	users.writes++
	users.lock.Unlock()
}

func GetCurrentLevelExperience(userId string) uint {
	users.lock.RLock()
	exp := users.cache[userId].CurrentLevelExperience
	users.lock.RUnlock()

	return exp
}

func UpdateAllUserInDb() error {
	for k, v := range users.cache {
		user := entities.UserType{}
		res := USERS_DB.FindOne(context.TODO(), bson.M{"complexid": v.ComplexID})
		res.Decode(&user)
		err := UpsertUserByID(k, &user)
		if err != nil {
			return err
		}
	}
	users.lock.Lock()
	users.writes = 0
	users.lock.Unlock()
	return nil
}

func GetWrites() int {
	users.lock.RLock()
	writes := users.writes
	users.lock.RUnlock()

	return writes
}

func StartUsersDBUpdater() {
	for {
		if GetWrites() > 5 {
			fmt.Println("updating db")
			UpdateAllUserInDb()
		}
	}
}

func RetrieveUsersByServerID(serverID string) []entities.UserType {
	var u []entities.UserType

	for _, v := range users {
		if v.ServerID == serverID {
			u = append(u, *v)
		}
	}

	return u
}
