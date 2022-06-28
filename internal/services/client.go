package services

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

	var perms int64
	var err error

	if perms, err = client.State.UserChannelPermissions(userId, serverId); err != nil {
	} else {
		perms, err = client.UserChannelPermissions(userId, serverId)
		if err != nil {
			log.Error("Error retrieving user permissions", err)
		}
	}

	if perms&discordgo.PermissionManageServer == 0 {
		return true
	}

	return false
}

func IsUserAdmin(userId, channelId string) bool {
	var perms int64
	var err error

	perms, err = client.State.UserChannelPermissions(userId, channelId)
	if err != nil {
		log.Error("Error retrieving user permissions", err)
	}

	if perms&discordgo.PermissionAdministrator == 0 {
		log.Info("User is not Admin")
		return false
	}

	log.Info("User is Admin")
	return true

}

func SetRole(userId, serverId, roleId string) error {

	err := client.GuildMemberRoleAdd(serverId, userId, roleId)
	if err != nil {
		return err
	}
	return nil
}

func RemoveRole(userId, serverId, roleId, oldRoleId string) error {

	err := client.GuildMemberRoleRemove(serverId, userId, oldRoleId)
	if err != nil {
		return err
	}
	return nil
}
