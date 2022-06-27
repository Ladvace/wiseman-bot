package main

import (
	"context"
	"fmt"
	"time"
	"wiseman/internal"
	"wiseman/internal/commands"
	"wiseman/internal/db"
	"wiseman/internal/services"

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
	d, err := services.Connect()
	if err != nil {
		panic(err)
	}
	defer d.Close()
	fmt.Println("Connected to Discord")

	// Hydrate data on cache
	start := time.Now()

	ns, err := db.HydrateServers(d)
	if err != nil {
		panic(err)
	}

	fmt.Println(ns, "Servers hydrated in", time.Since(start))
	start = time.Now()

	nu, err := db.HydrateUsers(d)
	if err != nil {
		panic(err)
	}

	fmt.Println(nu, "Users hydrated in", time.Since(start))

	db.HydrateCrossLookup()

	db.Hydrated = true

	services.StartHandlers()

	commands.Init()

	go db.StartUsersDBUpdater()

	// Start REST API
	e := internal.StartEcho()
	e.Logger.Fatal(e.Start(":1323"))

}
