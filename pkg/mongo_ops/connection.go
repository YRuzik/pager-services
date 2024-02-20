package mongo_ops

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
)

var Client *mongo.Client

type MongoCollections struct {
	TestCollection    *mongo.Collection
	ChatCollection    *mongo.Collection
	ProfileCollection *mongo.Collection
	UsersCollection   *mongo.Collection
	MembersCollection *mongo.Collection
}

var CollectionsPoll MongoCollections

func InitMongoDB() {
	DB_RS := os.Getenv("MDBREPLNAME")
	DB_NAME := os.Getenv("MDBNAME")
	DB_HOSTS := []string{
		"rc1a-7ei04i040frshb4c.mdb.yandexcloud.net:27018",
	}
	DB_USER := os.Getenv("MDBUSER")
	DB_PASS := os.Getenv("MDBPASSWORD")

	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?replicaSet=%s",
		DB_USER,
		DB_PASS,
		strings.Join(DB_HOSTS, ","),
		DB_NAME,
		DB_RS)

	if os.Getenv("DEBUG") == "LOCAL" {
		log.Print("start local")
		uri = "mongodb://localhost:27017"
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			panic(err)
		}
		Client = client
	} else {
		log.Print("start prods")
		//caCert, err := os.ReadFile("mdbcerts/root.crt")
		//if err != nil {
		//	panic(err)
		//}
		//caCertPool := x509.NewCertPool()
		//if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		//	panic("Error: CA file must be in PEM format")
		//}
		//tlsConfig := &tls.Config{
		//	RootCAs: caCertPool,
		//}
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
		Client = client
	}

	var result bson.M
	if err := Client.Database("local").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	database := Client.Database("db1")

	CollectionsPoll = MongoCollections{
		ChatCollection:    database.Collection("chats"),
		ProfileCollection: database.Collection("profiles"),
		UsersCollection:   database.Collection("profiles_md"),
		MembersCollection: database.Collection("members"),
	}
}
