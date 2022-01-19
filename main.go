package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"wiseman/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dg *discordgo.Session
var mongoClient *mongo.Client

type Server struct {
	ServerId            string   `bson:"serverid"`
	GuildPrefix         string   `bson:"guildprefix"`
	NotificationChannel string   `bson:"notificationchannel"`
	WelcomeChannel      string   `bson:"welcomechannel"`
	CustomRanks         []string `bson:"customranks"`
	RankTime            int      `bson:"ranktime"`
	WelcomeMessage      string   `bson:"welcomemessage"`
	DefaultRole         string   `bson:"defaultrole"`
}

type User struct {
	UserId        string `bson:"userid"`
	MessagesCount uint   `bson:"messagescount"`
	Rank          int    `bson:"rank"`
	Time          uint   `bson:"time"`
	Exp           uint   `bson:"exp"`
	GuildId       string `bson:"guildid"`
	LastRankTime  uint64 `bson:"lastranktime"`
	Bot           bool   `bson:"bot"`
	Verified      bool   `bson:"verified"`
}

var Servers map[string]Server
var Users map[string]User

func init() {
	Servers = make(map[string]Server)
	Users = make(map[string]User)
}

func setupDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := mongoClient.Database(utils.DB_NAME, nil)

	// Swallow errors
	db.CreateCollection(ctx, utils.USERS_INFIX)
	db.CreateCollection(ctx, utils.USERS_INFIX)

	return nil
}

func initServers() error {
	guilds, err := dg.UserGuilds(10, "", "")
	if err != nil {
		return err
	}

	// TODO: Use InsertMany to optimize this
	for _, guild := range guilds {
		// Check if server is already in DB
		res := mongoClient.Database(utils.DB_NAME).Collection(utils.SERVERS_INFIX).FindOne(context.TODO(), bson.M{"serverid": guild.ID})

		if res.Err() != mongo.ErrNoDocuments {
			var server Server
			err := res.Decode(&server)
			if err != nil {
				return err
			}
			Servers[guild.ID] = server
			continue
		}

		fmt.Println("Server not found in DB", guild.ID, guild.Name)
		server := Server{
			ServerId:            guild.ID,
			GuildPrefix:         "!",
			NotificationChannel: "",
			WelcomeChannel:      "",
			CustomRanks:         []string{},
			RankTime:            0,
			WelcomeMessage:      "",
			DefaultRole:         "",
		}
		Servers[guild.ID] = server

		mongoClient.Database(utils.DB_NAME).Collection(utils.SERVERS_INFIX).InsertOne(context.TODO(), server)
	}

	return nil
}

func initUsers() error {

	for k, _ := range Servers {

		members, err := dg.GuildMembers(k, "", 10)
		if err != nil {
			return err
		}

		// TODO: Use InsertMany to optimize this
		for _, member := range members {
			// Check if server is already in DB
			res := mongoClient.Database(utils.DB_NAME).Collection(utils.USERS_INFIX).FindOne(context.TODO(), bson.M{"userid": member.User.ID})
			if res.Err() != mongo.ErrNoDocuments {
				var user User
				err := res.Decode(&user)
				if err != nil {
					return err
				}
				Users[member.User.ID] = user
				continue
			}

			fmt.Println("User not found in DB", member.User.ID, member.User.Username+"#"+member.User.Discriminator)
			user := User{
				UserId:        member.User.ID,
				Verified:      member.User.Verified,
				Bot:           member.User.Bot,
				MessagesCount: 0,
				Rank:          0,
				Time:          0,
				Exp:           0,
				GuildId:       k,
				LastRankTime:  0,
			}

			mongoClient.Database(utils.DB_NAME).Collection(utils.USERS_INFIX).InsertOne(context.TODO(), user)

			Users[member.User.ID] = user
		}
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	err = setupDB()
	if err != nil {
		log.Fatal(err)
	}

	err = initServers()
	if err != nil {
		log.Fatal(err)
	}

	err = initUsers()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	defer dg.Close()
	defer cancel()

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
