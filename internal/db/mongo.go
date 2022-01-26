package db

import (
	"context"
	"os"
	"time"
	"wiseman/internal/shared"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var Hydrated bool

func Connect() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return nil, err
	}

	SERVERS_DB = mongoClient.Database(shared.DB_NAME).Collection(shared.SERVERS_INFIX)
	USERS_DB = mongoClient.Database(shared.DB_NAME).Collection(shared.USERS_INFIX)

	err = mongoClient.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return mongoClient, nil
}
