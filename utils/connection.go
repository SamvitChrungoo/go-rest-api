package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client -> db client
var Client *mongo.Client

//ConnectDB connects mongoDB
func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	password := GetEnvironmentVariable("DB_PASSWORD")
	Client, _ = mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://samvit:"+password+"@test.cmlur.mongodb.net/?retryWrites=true&w=majority"))
}

//DisconnectDB -> disconnects the database connection
func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Client.Disconnect(ctx)
}
