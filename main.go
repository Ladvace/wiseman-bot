package main

import (
	"context"
	"fmt"
	"time"
	"wiseman/internal"
	"wiseman/internal/commands"
	"wiseman/internal/db"
	"wiseman/internal/discord"

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

	// Initialize DB and collections
	db.SetupDB()
	fmt.Println("DB Initialized")

	// Connect to discord
	d, err := discord.Connect()
	if err != nil {
		panic(err)
	}
	defer d.Close()
	fmt.Println("Connected to Discord")

	// Hydrate data on cache
	start := time.Now()
	ns, err := db.HydrateServers(d)
	fmt.Println(ns, "Servers hydrated in", time.Since(start))
	start = time.Now()
	nu, err := db.HydrateUsers(d)
	fmt.Println(nu, "Users hydrated in", time.Since(start))

	discord.StartHandlers()

	commands.Init()

	// Start REST API
	e := internal.StartEcho()
	e.Logger.Fatal(e.Start(":1323"))

}
