package main

import (
	"assm/service-payment/db"
	"assm/service-payment/proto"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"time"
)

type Server struct {
	proto.UnimplementedPaymentServiceServer
	client     *mongo.Client
	collection *mongo.Collection
}

type KafkaEvent struct {
	Debtor   string
	Amount   float32
	Creditor string
	Event    string
}

func (s *Server) BalanceInquiry(ctx context.Context, req *proto.BalanceReq) (res *proto.BalanceRes, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {

			return &proto.BalanceRes{
				Result: &proto.BalanceRes_Failure{Failure: &proto.Failure{
					FailureCode:    proto.FailureCode_GENERAL_ERROR,
					FailureMessage: "User Id is empty",
				}},
			}, nil
		}

		existingBalance := db.GetBalance(s.collection, userId[0])

		return &proto.BalanceRes{
			Result: &proto.BalanceRes_Success_{
				Success: &proto.BalanceRes_Success{
					Balance: float32(existingBalance),
				},
			},
		}, nil
	}

	return &proto.BalanceRes{
		Result: &proto.BalanceRes_Failure{
			Failure: &proto.Failure{
				FailureCode:    proto.FailureCode_MISSING_DATA,
				FailureMessage: "User Id is empty",
			},
		},
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":9001")

	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()

	mongoClient := db.Conn()
	collection := db.Collc(mongoClient, "users")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{"halan"}, nil)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			msg, err := c.ReadMessage(time.Second)

			if err == nil {
				var event KafkaEvent
				err := bson.Unmarshal(msg.Value, &event)

				if err != nil {
					log.Fatal(err)
				}

				switch event.Event {
				case "Withdrawal":
					if event.Debtor == "" {
						log.Println("Invalid data")
						return
					}
					if !db.EditBalance(mongoClient, collection, event.Debtor, -event.Amount) {
						log.Println("failed to update balance")
						return
					}
				case "Deposit":
					if event.Amount < 0 || event.Creditor == "" {
						log.Println("Invalid data")
						return
					}
					if !db.EditBalance(mongoClient, collection, event.Creditor, event.Amount) {
						log.Println("failed to update balance")
						return
					}

				case "P2PTransfer":
					if event.Amount < 0 || event.Creditor == "" || event.Debtor == "" {
						log.Println("Invalid data")
						return
					}

					transferStatus, err := db.Transfer(mongoClient, collection, event.Debtor, event.Creditor, event.Amount)

					if err != nil {
						log.Println(err)
						return
					}

					if !transferStatus {
						log.Println("failed to update balance")
						return
					}
				}

			} else if !err.(kafka.Error).IsTimeout() {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}()

	proto.RegisterPaymentServiceServer(s, &Server{
		client:     mongoClient,
		collection: collection,
	})

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
