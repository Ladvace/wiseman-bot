package main

import (
	"context"
	"fmt"
	"wiseman/internal"
	"wiseman/internal/db"
	"wiseman/internal/discord"
	"wiseman/internal/servers"
	"wiseman/internal/users"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Connect to mongo
	mongo, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer mongo.Disconnect(context.TODO())

	// Connect to discord
	discord, err := discord.Connect()
	if err != nil {
		panic(err)
	}
	defer discord.Close()

	// Initialize DB and collections
	db.SetupDB()

	// Hydrate data on cache
	servers.Hydrate(discord, mongo)
	users.Hydrate(discord, mongo)

	// Start REST API
	e := internal.StartEcho()
	e.Logger.Fatal(e.Start(":1323"))

}
