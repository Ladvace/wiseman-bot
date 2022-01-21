package discord

import (
	"fmt"
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

	Client.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Register the messageCreate func as a callback for MessageCreate events.
	Client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, mongo)
	})

	Client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		fmt.Println(g.Name)
	})

	Client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		fmt.Println(g.ID)
	})

	// Open a websocket connection to Discord and begin listening.
	err = Client.Open()
	if err != nil {
		return nil, err
	}

	return Client, nil
}
