package db

import (
	"fmt"
	"wiseman/internal/entities"
)

var userServers = make(map[string][]*entities.ServerType, 50000)
var serverUsers = make(map[string][]*entities.UserType, 50000)

func HydrateCrossLookup() {
	for _, s := range users {
		// Check if user already has a servers list
		if _, ok := userServers[s.UserID]; !ok {
			userServers[s.UserID] = make([]*entities.ServerType, 0)
		}

		// Check if server already has a users list
		if _, ok := serverUsers[s.ServerID]; !ok {
			serverUsers[s.UserID] = make([]*entities.UserType, 0)
		}

		userServers[s.UserID] = append(userServers[s.UserID], servers[s.ServerID])
		serverUsers[s.ServerID] = append(serverUsers[s.ServerID], s)

	}
	fmt.Println("Hydrated cross lookup maps")
}

func GetServerUsers(userID string) []*entities.ServerType {
	return userServers[userID]
}
