package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Conn() *mongo.Client {
	_ = godotenv.Load()

	mongoUri := os.Getenv("MONGO_URI")
	collectionName := os.Getenv("MONGO_COLLECTION_NAME")

	if mongoUri == "" {
		panic("env variable MONGO_URI not found")
	}

	if collectionName == "" {
		panic("env variable MONGO_COLLECTION_NAME not found")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))

	if err != nil {
		fmt.Println("Couldn't connect to the database. \n verbose: ", err)
		panic("Couldn't connect to the database.")
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println("Database is offline. \n verbose: ", err)
		panic("Database is offline.")
	}

	return client
}

func Collc(collcName string) *mongo.Collection {
	conn := Conn()
	collec := conn.Database("halan").Collection(collcName)
	return collec
}
