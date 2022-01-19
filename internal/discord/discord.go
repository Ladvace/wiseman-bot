package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *discordgo.Session

func Connect(mongo *mongo.Client) (*discordgo.Session, error) {
	// Create a new Discord session using the provided bot token.
	var err error
	Client, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return nil, err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	Client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, mongo)
	})

	// In this example, we only care about receiving message events.
	Client.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = Client.Open()
	if err != nil {
		return nil, err
	}

	return Client, nil
}
