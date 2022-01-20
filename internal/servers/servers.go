package servers

import (
	"context"
	"fmt"
	"wiseman/internal/shared"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerType struct {
	ServerId            string            `bson:"serverid"`
	GuildPrefix         string            `bson:"guildprefix"`
	NotificationChannel string            `bson:"notificationchannel"`
	WelcomeChannel      string            `bson:"welcomechannel"`
	CustomRanks         map[string]string `bson:"customranks"`
	RankTime            int               `bson:"ranktime"`
	WelcomeMessage      string            `bson:"welcomemessage"`
	DefaultRole         string            `bson:"defaultrole"`
}

type ServersType map[string]ServerType

var servers ServersType

func init() {
	servers = make(map[string]ServerType, 1000)
}

func GetAll() *ServersType {
	return &servers
}

func Get(id string) ServerType {
	return servers[id]
}

func Upsert(id string, u ServerType) {
	servers[id] = u
}

func Hydrate(d *discordgo.Session, m *mongo.Client) error {
	var guilds []*discordgo.UserGuild
	var lastId string
	for {
		newGuilds, err := d.UserGuilds(100, "", lastId)
		if err != nil {
			return err
		}

		if len(newGuilds) == 0 {
			break
		}

		lastId = newGuilds[len(newGuilds)-1].ID
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
				return err
			}
			Upsert(guild.ID, server)
			continue
		}

		fmt.Println("Server not found in DB", guild.ID, guild.Name)
		server := ServerType{
			ServerId:            guild.ID,
			GuildPrefix:         "!",
			NotificationChannel: "",
			WelcomeChannel:      "",
			CustomRanks:         map[string]string{},
			RankTime:            0,
			WelcomeMessage:      "",
			DefaultRole:         "",
		}
		Upsert(guild.ID, server)

		m.Database(shared.DB_NAME).Collection(shared.SERVERS_INFIX).InsertOne(context.TODO(), server)
	}

	return nil
}
