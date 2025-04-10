package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func Conn() *mongo.Client {
	_ = godotenv.Load("../.env")

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

func FilterTodos(collection *mongo.Collection, query *bson.D) []bson.M {
	if query == nil {
		query = &bson.D{}
	}

	cursor, err := collection.Find(context.Background(), query)

	if err != nil {
		fmt.Println("Something went wrong: ", err)
		return nil
	}

	var results []bson.M

	if err = cursor.All(context.Background(), &results); err != nil {
		fmt.Println("Something went wrong: ", err)
		return nil
	}

	return results
}

func InsertTodo(collection *mongo.Collection, data any) bool {
	_, err := collection.InsertOne(context.Background(), &bson.D{
		{"data", data},
	})

	if err != nil {
		fmt.Println("[DB] Something went wrong: ", err)
		return false
	}

	return true
}

func DeleteTodo(collection *mongo.Collection, entryId primitive.ObjectID) bool {
	cursor, err := collection.DeleteOne(context.TODO(), bson.D{
		{"_id", entryId},
	})

	if err != nil {
		fmt.Println("Something went wrong: ", err)
		return false
	}

	if cursor.DeletedCount > 0 {
		return true
	}

	return false
}
