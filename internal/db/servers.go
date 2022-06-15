package db

import (
	"context"
	"fmt"
	"sort"
	"wiseman/internal/entities"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var servers entities.ServersType
var SERVERS_DB *mongo.Collection

func init() {
	servers = make(map[string]*entities.ServerType, 1000)
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
			var server entities.ServerType
			err := res.Decode(&server)
			if err != nil {
				return 0, err
			}

			// TODO: FIX
			sort.SliceStable(server.CustomRanks, func(i, j int) bool {
				return server.CustomRanks[i].MinLevel > server.CustomRanks[j].MinLevel
			})

			UpsertServerByID(guild.ID, &server)
			continue
		}

		fmt.Println("Server not found in DB", guild.ID, guild.Name)
		ns += 1

		server := entities.ServerType{
			ServerID:            guild.ID,
			ServerPrefix:        "!",
			NotificationChannel: "",
			WelcomeChannel:      "",
			CustomRanks:         []entities.RoleType{},
			RankTime:            0,
			MsgExpMultiplier:    1.00,
			TimeExpMultiplier:   1.00,
			WelcomeMessage:      "",
			DefaultRole:         "",
		}
		UpsertServerByID(guild.ID, &server)

		SERVERS_DB.InsertOne(context.TODO(), server)
	}

	return ns, nil
}

func GetServerByID(serverID string) *entities.ServerType {
	return servers[serverID]
}

func UpsertServerByID(serverID string, server *entities.ServerType) {
	servers[serverID] = server
}

func GetCustomRanksByGuildId(guildId string) []entities.RoleType {
	return servers[guildId].CustomRanks
}

func UpdateRoleServer(serverID string, rank entities.RoleType) error {

	servers[serverID].CustomRanks = append(servers[serverID].CustomRanks, rank)
	res := SERVERS_DB.FindOneAndUpdate(context.TODO(), bson.M{"serverid": serverID}, bson.M{"$set": bson.M{"customranks": servers[serverID].CustomRanks}})

	return res.Err()
}

func GetRankRoleByLevel(s entities.ServerType, level uint) entities.RoleType {
	for _, v := range s.CustomRanks {
		if level >= v.MinLevel {
			return v
		}
	}

	return entities.RoleType{
		Id:       "",
		MinLevel: 0,
	}
}

func GetServerMultiplierByGuildId(guildId string) float64 {
	return servers[guildId].MsgExpMultiplier
}
