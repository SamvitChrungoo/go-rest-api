package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ConnectDB connects mongoDB
func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	password := GetEnvironmentVariable("DB_PASSWORD")
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://samvit:"+password+"@test.cmlur.mongodb.net/?retryWrites=true&w=majority"))
	return client
}
