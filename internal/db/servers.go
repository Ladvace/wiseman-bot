package db

import (
	"context"
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RankType struct {
	RankName     string `bson:"rankname"`
	RankMinLevel uint   `bson:"rankminlevel"`
}

type ServerType struct {
	ServerID            string     `bson:"serverid"`
	ServerPrefix        string     `bson:"guildprefix"`
	NotificationChannel string     `bson:"notificationchannel"`
	WelcomeChannel      string     `bson:"welcomechannel"`
	CustomRanks         []RankType `bson:"customranks"`
	RankTime            int        `bson:"ranktime"`
	MsgExpMultiplier    float64    `bson:"msgexpmultiplier"`
	TimeExpMultiplier   float64    `bson:"timeexpmultiplier"`
	WelcomeMessage      string     `bson:"welcomemessage"`
	DefaultRole         string     `bson:"defaultrole"`
}

type ServersType map[string]ServerType

var servers ServersType

var SERVERS_DB *mongo.Collection

func init() {
	servers = make(map[string]ServerType, 1000)
}

func (s ServerType) GetRankRoleByLevel(level uint) RankType {
	for _, v := range s.CustomRanks {
		if level >= v.RankMinLevel {
			return v
		}
	}

	return RankType{
		RankName:     "",
		RankMinLevel: 0,
	}
}

func HydrateServers(d *discordgo.Session) (int, error) {
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
		res := SERVERS_DB.FindOne(context.TODO(), bson.M{"serverid": guild.ID})

		if res.Err() != mongo.ErrNoDocuments {
			var server ServerType
			err := res.Decode(&server)
			if err != nil {
				return 0, err
			}

			// TODO: FIX
			sort.SliceStable(server.CustomRanks, func(i, j int) bool {
				return server.CustomRanks[i].RankMinLevel > server.CustomRanks[j].RankMinLevel
			})

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
			CustomRanks:         []RankType{},
			RankTime:            0,
			MsgExpMultiplier:    1.00,
			TimeExpMultiplier:   1.00,
			WelcomeMessage:      "",
			DefaultRole:         "",
		}
		UpsertServerByID(guild.ID, server)

		SERVERS_DB.InsertOne(context.TODO(), server)
	}

	return ns, nil
}

func GetServerByID(serverID string) ServerType {
	return servers[serverID]
}

func UpsertServerByID(serverID string, server ServerType) {
	servers[serverID] = server
}
