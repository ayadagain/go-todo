package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
)

type TodoModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Data      string             `bson:"data"`
	CreatedBy string             `bson:"created_by"`
}

type UserModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	balance  float64            `bson:"balance"`
	username string             `bson:"username"`
}

type PaymentsModel struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	from   primitive.ObjectID `bson:"from"`
	to     primitive.ObjectID `bson:"to"`
	amount float64            `bson:"amount"`
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

func Collc(client *mongo.Client, collcName string) *mongo.Collection {
	collec := client.Database("halan").Collection(collcName)
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

func GetBalance(collection *mongo.Collection, userId string) float64 {
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatalf("Failed to convert string to ObjectID: %v", err)
	}

	cursor := collection.FindOne(context.Background(), bson.D{
		{"_id", objectID},
	})

	var result bson.M

	if err := cursor.Decode(&result); err != nil {
		log.Fatalf("Failed to convert string to ObjectID: %v", err)
	}

	return result["balance"].(float64)
}

func EditBalance(client *mongo.Client, collection *mongo.Collection, userId string, amount float32) bool {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatalf("Failed to convert string to ObjectID: %v", err)
		return false
	}

	session, err := client.StartSession()
	defer session.EndSession(context.TODO())

	if err != nil {
		log.Fatalf("Failed to start session: %v", err)
		return false
	}

	result, err := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		result, err := collection.UpdateOne(
			ctx,
			bson.D{{
				"_id", objectID,
			}},
			bson.D{
				{"$inc", bson.D{
					{"balance", amount},
				}},
			},
		)
		return *result, err
	}, txnOptions)

	if err != nil {
		log.Fatalf("Failed to update: %v", err)
		return false
	}

	return true
}

func Transfer(client *mongo.Client, userCollection *mongo.Collection, from string, to string, amount float32) (bool, error) {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	fromUserId, err := primitive.ObjectIDFromHex(from)
	if err != nil {
		return false, status.Errorf(codes.InvalidArgument, "Failed to convert string to ObjectID: %v", err)
	}

	toUserId, err := primitive.ObjectIDFromHex(to)
	if err != nil {
		return false, status.Errorf(codes.InvalidArgument, "Failed to convert string to ObjectID: %v", err)
	}

	if amount < 1 {
		return false, status.Errorf(codes.InvalidArgument, "amount must be greater than zero")
	}

	fromCursor := userCollection.FindOne(context.Background(), bson.D{
		{"_id", fromUserId},
	})

	if err := fromCursor.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, status.Errorf(codes.NotFound, "User %s not found", fromUserId)
		}
		return false, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	toCursor := userCollection.FindOne(context.Background(), bson.D{
		{"_id", toUserId},
	})

	if err := toCursor.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, status.Errorf(codes.NotFound, "User %s not found", toUserId)
		}
		return false, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	var fromUser bson.M
	var toUser bson.M

	if err := fromCursor.Decode(&fromUser); err != nil {
		return false, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	if err := toCursor.Decode(&toUser); err != nil {
		return false, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	fromBalance := float32(fromUser["balance"].(float64))

	if fromBalance < amount {
		return false, status.Errorf(codes.OutOfRange, "Amount insufficient")
	}

	session, err := client.StartSession()
	if err != nil {
		return false, status.Errorf(codes.Internal, "Database error: %v", err)
	}
	defer session.EndSession(context.TODO())
	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		fromResult, fromErr := userCollection.UpdateOne(
			ctx,
			bson.D{{
				"_id", fromUserId,
			}},
			bson.D{
				{"$inc", bson.D{
					{"balance", -amount},
				}},
			},
		)

		if fromErr != nil {
			return nil, fromErr
		}

		if fromResult.MatchedCount == 0 {
			return nil, status.Errorf(codes.NotFound, "['from'] User %s not found", fromUserId)
		}

		toResult, toErr := userCollection.UpdateOne(
			ctx,
			bson.D{{
				"_id", toUserId,
			}},
			bson.D{
				{"$inc", bson.D{
					{"balance", amount},
				}},
			},
		)

		if toErr != nil {
			return nil, toErr
		}
		if toResult.MatchedCount == 0 {
			return nil, status.Errorf(codes.NotFound, "User %s not found", fromUserId)
		}

		return true, nil
	}, txnOptions)

	if err != nil {
		fmt.Println("Transfer err ", err)
	}

	return true, nil
}
