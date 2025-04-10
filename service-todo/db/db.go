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

type TodoModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Data      string             `bson:"data" json:"k"`
	CreatedBy string             `bson:"created_by" json:"created_by"`
}

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

func FilterTodos(collection *mongo.Collection, query *bson.D) []TodoModel {
	if query == nil {
		query = &bson.D{}
	}

	cursor, err := collection.Find(context.Background(), query)

	if err != nil {
		fmt.Println("Something went wrong: ", err)
		return nil
	}

	var result []TodoModel

	if err = cursor.All(context.Background(), &result); err != nil {
		fmt.Println("Something went wrong: ", err)
		return nil
	}

	return result
}

func InsertTodo(collection *mongo.Collection, data any, userId string) bool {
	_, err := collection.InsertOne(context.Background(), &TodoModel{
		Data:      data.(string),
		CreatedBy: userId,
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
