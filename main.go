package main

import (
	"context"
	"fmt"
	"wiseman/internal"
	"wiseman/internal/commands"
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
	fmt.Println("Connected to MongoDB")

	// Connect to discord
	discord, err := discord.Connect(mongo)
	if err != nil {
		panic(err)
	}
	defer discord.Close()
	fmt.Println("Connected to Discord")

	// Initialize DB and collections
	db.SetupDB()
	fmt.Println("DB Initialized")

	// Hydrate data on cache
	servers.Hydrate(discord, mongo)
	fmt.Println("Servers hydrated")
	users.Hydrate(discord, mongo)
	fmt.Println("Users hydrated")

	commands.Init()

	// Start REST API
	e := internal.StartEcho()
	e.Logger.Fatal(e.Start(":1323"))

}
