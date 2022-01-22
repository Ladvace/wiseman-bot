package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var client *discordgo.Session

func Connect() (*discordgo.Session, error) {
	// Create a new Discord session using the provided bot token.
	var err error
	client, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return nil, err
	}

	client.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Open a websocket connection to Discord and begin listening.
	err = client.Open()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func StartHandlers() {
	// Register the messageCreate func as a callback for MessageCreate events.
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m)
	})

	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		serverAdd(s, g)
	})

	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		serverRemove(s, g)
	})

	client.AddHandler(func(s *discordgo.Session, u *discordgo.GuildMemberAdd) {
		memberAdd(s, u)
	})

	client.AddHandler(func(s *discordgo.Session, u *discordgo.GuildMemberRemove) {
		memberRemove(s, u)
	})
}
