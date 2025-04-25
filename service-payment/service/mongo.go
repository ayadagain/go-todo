package service

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
)

var (
	ConversionErr = errors.New("conversion error")
	DecodeErr     = errors.New("decode error")
	SessionErr    = errors.New("session error")
	QueryErr      = errors.New("query error")
)

func (s DefaultPaymentService) GetBalance(userId string) (float64, error) {
	objectID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return -1, ConversionErr
	}

	mongoCollection := s.ServiceContext.MongoCollection()

	cursor := mongoCollection.FindOne(context.Background(), bson.M{"_id": objectID})

	var result bson.M

	if err := cursor.Decode(&result); err != nil {
		return -1, DecodeErr
	}

	return result["balance"].(float64), nil
}

func (s DefaultPaymentService) EditBalance(userId string, amount float64) (bool, error) {
	txnOptions := options.Transaction().SetWriteConcern(writeconcern.Majority())

	objectID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return false, ConversionErr
	}

	mongoClient := s.ServiceContext.MongoClient()
	mongoCollection := s.ServiceContext.MongoCollection()
	session, err := mongoClient.StartSession()

	defer session.EndSession(context.Background())

	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return false, SessionErr
	}

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		result, err := mongoCollection.UpdateOne(
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
		log.Fatalf("Failed to update balance: %v", err)
		return false, QueryErr
	}

	return true, nil
}
