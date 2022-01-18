package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	var dg *discordgo.Session
	defer dg.Close()

	go func() {
		// Create a new Discord session using the provided bot token.
		dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
		if err != nil {
			fmt.Println("error creating Discord session,", err)
			return
		}

		// Register the messageCreate func as a callback for MessageCreate events.
		dg.AddHandler(messageCreate)

		// In this example, we only care about receiving message events.
		dg.Identify.Intents = discordgo.IntentsGuildMessages

		// Open a websocket connection to Discord and begin listening.
		err = dg.Open()
		if err != nil {
			fmt.Println("error opening connection,", err)
			return
		}

		// Wait here until CTRL-C or other term signal is received.
		fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	}()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
