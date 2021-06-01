package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

//ConnectDB connects mongoDB
func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	password := "samvit123"
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://samvit:"+password+"@test.cmlur.mongodb.net/?retryWrites=true&w=majority"))
	return client
}

//DisconnectDB disconnects mongoDB
func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Disconnect(ctx)
}
