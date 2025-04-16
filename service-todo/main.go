package main

import (
	paymentProto "assm/service-payment/proto"
	hgrpc "assm/service-todo/grpc"
	"assm/service-todo/proto"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedTodoServiceServer
	kafka      *kafka.Producer
	kTopic     string
	grpcClient paymentProto.PaymentServiceClient
}

func (s *Server) Withdraw(ctx context.Context, req *proto.WithdrawReq) (res *proto.WithdrawRes, err error) {
	if req.Amount == 0 {
		return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
			FailureCode:    proto.T_FailureCode_T_MISSING_DATA,
			FailureMessage: "Amount is empty",
		}}}, nil
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_MISSING_DATA,
				FailureMessage: "userId is empty",
			}}}, nil
		}

		md := metadata.New(map[string]string{"userId": userId[0]})
		ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

		existingBalance, err := s.grpcClient.BalanceInquiry(ctxWithMd, &paymentProto.BalanceReq{})

		if err != nil {
			return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_NETWORK_ERROR,
				FailureMessage: "Error calling the rpc",
			}}}, nil
		}

		switch result := existingBalance.Result.(type) {
		case *paymentProto.BalanceRes_Success_:
			if result.Success.Balance < req.Amount {
				return &proto.WithdrawRes{
					Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
						FailureCode:    proto.T_FailureCode_T_INSUFFICIENT_BALANCE,
						FailureMessage: "Insufficient funds",
					}},
				}, nil
			}

			kMessage := bson.D{
				{Key: "debtor", Value: userId[0]},
				{Key: "amount", Value: req.Amount},
				{Key: "creditor", Value: nil},
				{Key: "event", Value: "Withdrawal"},
			}

			kMessageBytes, err := bson.Marshal(kMessage)
			if err != nil {
				return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
					FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
					FailureMessage: "Something went wrong",
				}}}, nil
			}

			err = s.kafka.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &s.kTopic,
					Partition: kafka.PartitionAny,
				},
				Value: kMessageBytes,
			}, nil)

			if err != nil {
				return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
					FailureCode:    proto.T_FailureCode_T_NETWORK_ERROR,
					FailureMessage: "Kafka network error",
				}}}, nil
			}

			return &proto.WithdrawRes{
				Result: &proto.WithdrawRes_Success_{Success: &proto.WithdrawRes_Success{
					Status:  1,
					Message: fmt.Sprintf("You withdrew $%.2f from your balance. Your balance now is $%.2f", req.Amount, result.Success.Balance-req.Amount),
				}},
			}, nil

		case *paymentProto.BalanceRes_Failure:
			if result.Failure.FailureCode == paymentProto.FailureCode_MISSING_DATA {
				return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
					FailureCode:    proto.T_FailureCode_T_MISSING_DATA,
					FailureMessage: "Missing data",
				}}}, nil

			} else if result.Failure.FailureCode == paymentProto.FailureCode_INSUFFICIENT_BALANCE {
				return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
					FailureCode:    proto.T_FailureCode_T_INSUFFICIENT_BALANCE,
					FailureMessage: "Insufficient Balance",
				}}}, nil
			}

			return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
				FailureMessage: "Something went wrong",
			}}}, nil
		}
	}

	return &proto.WithdrawRes{Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
		FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
		FailureMessage: "Something went wrong",
	}}}, nil
}

func (s *Server) Deposit(ctx context.Context, req *proto.DepositReq) (res *proto.DepositRes, err error) {
	if req.Amount == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "amount is empty")
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return &proto.DepositRes{Result: &proto.DepositRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_MISSING_DATA,
				FailureMessage: "userId is empty",
			}}}, nil
		}

		kMessage := bson.D{
			{Key: "debtor", Value: nil},
			{Key: "amount", Value: req.Amount},
			{Key: "creditor", Value: userId[0]},
			{Key: "event", Value: "Deposit"},
		}

		kMessageBytes, err := bson.Marshal(kMessage)
		if err != nil {
			return &proto.DepositRes{Result: &proto.DepositRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
				FailureMessage: "Something went wrong",
			}}}, nil
		}

		err = s.kafka.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &s.kTopic,
				Partition: kafka.PartitionAny,
			},
			Value: kMessageBytes,
		}, nil)

		if err != nil {
			return &proto.DepositRes{Result: &proto.DepositRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_NETWORK_ERROR,
				FailureMessage: "Kafka network error",
			}}}, nil
		}

		return &proto.DepositRes{Result: &proto.DepositRes_Success_{Success: &proto.DepositRes_Success{
			Status:  1,
			Message: fmt.Sprintf("You deposited $%.2f in your balance", req.Amount),
		}}}, nil
	}

	return &proto.DepositRes{Result: &proto.DepositRes_Failure{Failure: &proto.T_Failure{
		FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
		FailureMessage: "Something went wrong",
	}}}, nil
}

func (s *Server) Transfer(ctx context.Context, req *proto.TransferReq) (res *proto.TransferRes, err error) {
	if req.Amount == 0 || req.To == "" {
		return nil, status.Errorf(codes.InvalidArgument, "incomplete request")
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return &proto.TransferRes{Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_MISSING_DATA,
				FailureMessage: "userId is empty",
			}}}, nil
		}

		md := metadata.New(map[string]string{"userId": userId[0]})
		ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

		existingBalance, err := s.grpcClient.BalanceInquiry(ctxWithMd, &paymentProto.BalanceReq{})

		if err != nil {
			return &proto.TransferRes{Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_NETWORK_ERROR,
				FailureMessage: "Error calling the rpc",
			}}}, nil
		}

		switch result := existingBalance.Result.(type) {
		case *paymentProto.BalanceRes_Success_:
			if result.Success.Balance < req.Amount {
				return &proto.TransferRes{
					Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
						FailureCode:    proto.T_FailureCode_T_INSUFFICIENT_BALANCE,
						FailureMessage: "Insufficient funds",
					}},
				}, nil
			}

			kMessage := bson.D{
				{Key: "debtor", Value: userId[0]},
				{Key: "amount", Value: req.Amount},
				{Key: "creditor", Value: req.To},
				{Key: "event", Value: "P2PTransfer"},
			}

			kMessageBytes, err := bson.Marshal(kMessage)
			if err != nil {
				fmt.Println("damn bro: ", err.Error())
				return &proto.TransferRes{Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
					FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
					FailureMessage: "Something went wrong",
				}}}, nil
			}

			err = s.kafka.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &s.kTopic,
					Partition: kafka.PartitionAny,
				},
				Value: kMessageBytes,
			}, nil)

			return &proto.TransferRes{Result: &proto.TransferRes_Success_{Success: &proto.TransferRes_Success{
				Status:  1,
				Message: fmt.Sprintf("You have transferred $%.2f successfully.", req.Amount),
			}}}, nil

		case *paymentProto.BalanceRes_Failure:

			fmt.Println("testing.....: ", result.Failure)

			if result.Failure.FailureCode == paymentProto.FailureCode_MISSING_DATA {
				return &proto.TransferRes{Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
					FailureCode:    proto.T_FailureCode_T_MISSING_DATA,
					FailureMessage: "Missing data",
				}}}, nil

			} else if result.Failure.FailureCode == paymentProto.FailureCode_INSUFFICIENT_BALANCE {
				return &proto.TransferRes{Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
					FailureCode:    proto.T_FailureCode_T_INSUFFICIENT_BALANCE,
					FailureMessage: "Insufficient Balance",
				}}}, nil
			}

			return &proto.TransferRes{Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
				FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
				FailureMessage: "Something went wrong",
			}}}, nil
		}

	}
	return &proto.TransferRes{Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
		FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
		FailureMessage: "Something went wrong",
	}}}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	if err != nil {
		os.Exit(1)
	}

	defer kafkaProducer.Close()

	go func() {
		for e := range kafkaProducer.Events() {
			switch ev := (e).(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	proto.RegisterTodoServiceServer(s, &Server{
		kafka:      kafkaProducer,
		kTopic:     "halan",
		grpcClient: hgrpc.GrpcConn(),
	})

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
