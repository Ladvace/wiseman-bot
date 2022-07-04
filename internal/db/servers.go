package db

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"wiseman/internal/entities"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// var servers entities.ServersType
var SERVERS_DB *mongo.Collection

type Servers struct {
	cache  map[string]*entities.ServerType
	writes int
	lock   sync.RWMutex
}

var servers Servers = Servers{
	cache:  make(map[string]*entities.ServerType, 1000),
	writes: 0,
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
			CustomRanks:         []entities.CustomRanks{},
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
	servers.lock.RLock()
	s := servers.cache[serverID]
	servers.lock.RUnlock()

	return s
}

func UpsertServerByID(serverID string, server *entities.ServerType) {
	servers.lock.Lock()
	servers.cache[serverID] = server
	servers.writes++
	servers.lock.Unlock()
}

func GetCustomRanksByGuildId(guildId string) []entities.RoleType {
	servers.lock.Lock()
	cr := servers.cache[guildId].CustomRanks
	servers.lock.Unlock()

	return cr
}

func UpdateRoleServer(serverID string, rank entities.RoleType) {

	servers.lock.Lock()
	servers.cache[serverID].CustomRanks = append(servers.cache[serverID].CustomRanks, rank)
	servers.writes++
	// res := SERVERS_DB.FindOneAndUpdate(context.TODO(), bson.M{"serverid": serverID}, bson.M{"$set": bson.M{"customranks": servers.cache[serverID].CustomRanks}})
	servers.lock.Unlock()
	// return res.Err()
}

func GetRankRoleByLevel(s entities.ServerType, level uint) entities.CustomRanks {
	for _, v := range s.CustomRanks {
		if level >= v.MinLevel {
			return v
		}
	}

	return entities.CustomRanks{
		Id:       "",
		MinLevel: 0,
	}
}

func GetServerMultiplierByGuildId(guildId string) float64 {
	servers.lock.RLock()
	mem := servers.cache[guildId].MsgExpMultiplier
	servers.lock.RUnlock()

	return mem

}

func GetServersWrites() int {
	users.lock.RLock()
	writes := users.writes
	users.lock.RUnlock()

	return writes
}

func StartServersDBUpdater() {
	for {
		if GetWrites() > 5 {
			fmt.Println("updating server db")
			UpdateAllServersInDb()
		}
	}
}

func UpdateAllServersInDb() error {
	for k, v := range servers.cache {
		server := entities.ServerType{}
		res := USERS_DB.FindOne(context.TODO(), bson.M{"complexid": v})
		res.Decode(&server)
		UpsertServerByID(k, &server)
	}
	users.lock.Lock()
	users.writes = 0
	users.lock.Unlock()
	return nil
}
