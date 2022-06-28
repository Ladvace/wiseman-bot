package services

import "github.com/bwmarrin/discordgo"

func RetrieveServer(serverID string) (*discordgo.Guild, error) {
	g, err := client.State.Guild(serverID)

	if err == nil && g.ID != "" {
		return g, nil
	}

	g, err = client.Guild(serverID)
	if err == nil && g.ID != "" {
		return g, nil
	}

	return nil, err
}
