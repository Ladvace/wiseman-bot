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

func Connect() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return nil, err
	}

	err = mongoClient.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return mongoClient, nil
}

func SetupDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := mongoClient.Database(shared.DB_NAME, nil)

	// Swallow errors
	db.CreateCollection(ctx, shared.SERVERS_INFIX)
	db.CreateCollection(ctx, shared.USERS_INFIX)
}
