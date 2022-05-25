package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/gommon/log"
)

func RetrieveUser(userID, serverID string) (*discordgo.User, error) {
	u, err := client.State.Member(userID, serverID)

	if err == nil {
		return u.User, nil
	}

	user, err := client.User(userID)
	if err == nil {
		return user, nil
	}

	return nil, err
}

func IsUserManager(userId, serverId string) bool {

	perms, err := client.State.UserChannelPermissions(userId, serverId)
	if err != nil {
		log.Error("Error retrieving user permissions", err)
	}

	if perms&discordgo.PermissionManageServer == 0 {
		return true
	}

	return false
}

func IsUserAdmin(userId, serverId string) bool {
	perms, err := client.State.UserChannelPermissions(userId, serverId)
	if err != nil {
		log.Error("Error retrieving user permissions", err)
	}

	if perms&discordgo.PermissionAdministrator == 0 {
		return true
	}

	return false
}
