package database

import (
	"context"
	"log"

	"gariwallet/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.Global.MongoDBUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Cannot connect DB")
	}
	log.Println("Success to connect DB")
	return client
}

func DisconnectDB(ctx context.Context, mongoClient *mongo.Client) {
	err := mongoClient.Disconnect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connection to MongoDB closed.")
}
