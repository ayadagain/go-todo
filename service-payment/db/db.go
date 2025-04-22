package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

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

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
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
