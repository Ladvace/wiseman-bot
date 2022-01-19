package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var Client *discordgo.Session

func Connect() (*discordgo.Session, error) {
	// Create a new Discord session using the provided bot token.
	var err error
	Client, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return nil, err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	Client.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	Client.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = Client.Open()
	if err != nil {
		return nil, err
	}

	return Client, nil
}
