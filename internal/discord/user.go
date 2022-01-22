package discord

import "github.com/bwmarrin/discordgo"

func RetrieveUser(userID, serverID string) (*discordgo.User, error) {
	u, err := client.State.Member(userID, serverID)

	if err == nil && u.User.ID != "" {
		return u.User, nil
	}

	user, err := client.User(userID)
	if err == nil && user.ID != "" {
		return user, nil
	}

	return nil, err
}
