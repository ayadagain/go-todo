package service

import (
	paymentProto "assm/service-payment/proto"
	"assm/service-todo/client"
	"assm/service-todo/ctx"
	"assm/service-todo/proto"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type DefaultBankingService struct {
	ServiceContext ctx.ServiceCtx
	grpcClient     *client.PaymentClient
}

func NewDefaultBankingService(serviceContext ctx.ServiceCtx) *DefaultBankingService {
	grpcClient := client.NewPaymentClient(serviceContext.GrpcClient())
	return &DefaultBankingService{
		ServiceContext: serviceContext,
		grpcClient:     grpcClient,
	}
}

func (s DefaultBankingService) Withdraw(ctx context.Context, req *proto.WithdrawReq) (res *proto.WithdrawRes, err error) {
	kafkaP := s.ServiceContext.KafkaProducer()
	kafkaTopic := s.ServiceContext.Config().GetKafkaTopic()

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

		existingBalance, err := s.grpcClient.Client.BalanceInquiry(ctxWithMd, &paymentProto.BalanceReq{})

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

			err = kafkaP.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &kafkaTopic,
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

func (s DefaultBankingService) Deposit(ctx context.Context, req *proto.DepositReq) (res *proto.DepositRes, err error) {
	if req.Amount == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "amount is empty")
	}

	kafkaP := s.ServiceContext.KafkaProducer()
	kafkaTopic := s.ServiceContext.Config().GetKafkaTopic()

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

		err = kafkaP.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &kafkaTopic,
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

func (s DefaultBankingService) Transfer(ctx context.Context, req *proto.TransferReq) (res *proto.TransferRes, err error) {
	if req.Amount == 0 || req.To == "" {
		return nil, status.Errorf(codes.InvalidArgument, "incomplete request")
	}

	kafkaP := s.ServiceContext.KafkaProducer()
	kafkaTopic := s.ServiceContext.Config().GetKafkaTopic()

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

		existingBalance, err := s.grpcClient.Client.BalanceInquiry(ctxWithMd, &paymentProto.BalanceReq{})

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
				return &proto.TransferRes{Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
					FailureCode:    proto.T_FailureCode_T_GENERAL_ERROR,
					FailureMessage: "Something went wrong",
				}}}, nil
			}

			err = kafkaP.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &kafkaTopic,
					Partition: kafka.PartitionAny,
				},
				Value: kMessageBytes,
			}, nil)

			return &proto.TransferRes{Result: &proto.TransferRes_Success_{Success: &proto.TransferRes_Success{
				Status:  1,
				Message: fmt.Sprintf("You have transferred $%.2f successfully.", req.Amount),
			}}}, nil

		case *paymentProto.BalanceRes_Failure:
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
