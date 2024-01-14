package mongo_ops

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

type MongoCollections struct {
	TestCollection    *mongo.Collection
	ChatCollection    *mongo.Collection
	ProfileCollection *mongo.Collection
}

var CollectionsPoll MongoCollections

func InitMongoDB() {
	uri := "mongodb://localhost:27017"

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	var result bson.M
	if err := client.Database("local").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	Client = client
	database := Client.Database("test_streams")

	CollectionsPoll = MongoCollections{
		TestCollection:    database.Collection("transfers"),
		ChatCollection:    database.Collection("chats"),
		ProfileCollection: database.Collection("profiles"),
	}
}
