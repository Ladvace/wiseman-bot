package db

import (
	"context"
	"fmt"
	"wiseman/internal/shared"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerType struct {
	ServerID            string            `bson:"serverid"`
	ServerPrefix        string            `bson:"guildprefix"`
	NotificationChannel string            `bson:"notificationchannel"`
	WelcomeChannel      string            `bson:"welcomechannel"`
	CustomRanks         map[string]string `bson:"customranks"`
	RankTime            int               `bson:"ranktime"`
	WelcomeMessage      string            `bson:"welcomemessage"`
	DefaultRole         string            `bson:"defaultrole"`
}

type ServersType map[string]ServerType

var users UsersType
var servers ServersType

func init() {
	servers = make(map[string]ServerType, 1000)
}

func HydrateServers(d *discordgo.Session, m *mongo.Client) (int, error) {
	var ns int
	var guilds []*discordgo.UserGuild
	var lastID string
	for {
		newGuilds, err := d.UserGuilds(100, "", lastID)
		if err != nil {
			return 0, err
		}

		if len(newGuilds) == 0 {
			break
		}

		lastID = newGuilds[len(newGuilds)-1].ID
		guilds = append(guilds, newGuilds...)
	}

	// TODO: Use InsertMany to optimize this
	for _, guild := range guilds {
		// Check if server is already in DB
		res := m.Database(shared.DB_NAME).Collection(shared.SERVERS_INFIX).FindOne(context.TODO(), bson.M{"serverid": guild.ID})

		if res.Err() != mongo.ErrNoDocuments {
			var server ServerType
			err := res.Decode(&server)
			if err != nil {
				return 0, err
			}
			UpsertServerByID(guild.ID, server)
			continue
		}

		fmt.Println("Server not found in DB", guild.ID, guild.Name)
		ns += 1

		server := ServerType{
			ServerID:            guild.ID,
			ServerPrefix:        "!",
			NotificationChannel: "",
			WelcomeChannel:      "",
			CustomRanks:         map[string]string{},
			RankTime:            0,
			WelcomeMessage:      "",
			DefaultRole:         "",
		}
		UpsertServerByID(guild.ID, server)

		m.Database(shared.DB_NAME).Collection(shared.SERVERS_INFIX).InsertOne(context.TODO(), server)
	}

	return ns, nil
}

func GetServerByID(serverID string) ServerType {
	return servers[serverID]
}

func UpsertServerByID(serverID string, server ServerType) {
	servers[serverID] = server
}
