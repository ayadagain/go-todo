package service

import (
	paymentProto "assm/service-payment/proto"
	"assm/service-todo/client"
	"assm/service-todo/ctx"
	"assm/service-todo/proto"
	"assm/service-todo/response"
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
	kafkaProducer  *kafka.Producer
	kafkaTopic     string
}

func NewDefaultBankingService(serviceContext ctx.ServiceCtx) *DefaultBankingService {
	return &DefaultBankingService{
		ServiceContext: serviceContext,
		grpcClient:     client.NewPaymentClient(serviceContext.GrpcClient()),
		kafkaProducer:  serviceContext.KafkaProducer(),
		kafkaTopic:     serviceContext.Config().GetKafkaTopic(),
	}
}

func (s DefaultBankingService) Withdraw(ctx context.Context, req *proto.WithdrawReq) (res *proto.WithdrawRes, err error) {
	if req.Amount == 0 {
		return response.WithdrawResFailure(proto.T_FailureCode_T_MISSING_DATA, "Amount is empty"), nil
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return response.WithdrawResFailure(proto.T_FailureCode_T_MISSING_DATA, "User Id is empty"), nil
		}

		md := metadata.New(map[string]string{"userId": userId[0]})
		ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

		existingBalance, err := s.grpcClient.Client.BalanceInquiry(ctxWithMd, &paymentProto.BalanceReq{})

		if err != nil {
			return response.WithdrawResFailure(proto.T_FailureCode_T_NETWORK_ERROR, "Error calling the rpc"), nil
		}

		switch result := existingBalance.Result.(type) {
		case *paymentProto.BalanceRes_Success_:
			if result.Success.Balance < req.Amount {
				return response.WithdrawResFailure(proto.T_FailureCode_T_INSUFFICIENT_BALANCE, "Insufficient funds"), nil
			}

			kMessage := bson.D{
				{Key: "debtor", Value: userId[0]},
				{Key: "amount", Value: req.Amount},
				{Key: "creditor", Value: nil},
				{Key: "event", Value: "Withdrawal"},
			}

			kMessageBytes, err := bson.Marshal(kMessage)
			if err != nil {
				return response.WithdrawResFailure(proto.T_FailureCode_T_GENERAL_ERROR, "Something went wrong"), nil
			}

			err = s.kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &s.kafkaTopic,
					Partition: kafka.PartitionAny,
				},
				Value: kMessageBytes,
			}, nil)

			if err != nil {
				return response.WithdrawResFailure(proto.T_FailureCode_T_NETWORK_ERROR, "Kafka network error"), nil
			}

			return response.WithdrawResSuccess(fmt.Sprintf("You withdrew $%.2f from your balance. Your balance now is $%.2f", req.Amount, result.Success.Balance-req.Amount)), nil

		case *paymentProto.BalanceRes_Failure:
			if result.Failure.FailureCode == paymentProto.FailureCode_MISSING_DATA {
				return response.WithdrawResFailure(proto.T_FailureCode_T_MISSING_DATA, "Missing data"), nil
			} else if result.Failure.FailureCode == paymentProto.FailureCode_INSUFFICIENT_BALANCE {
				return response.WithdrawResFailure(proto.T_FailureCode_T_INSUFFICIENT_BALANCE, "Insufficient Balance"), nil
			}
			return response.WithdrawResFailure(proto.T_FailureCode_T_GENERAL_ERROR, "Something went wrong"), nil
		}
	}

	return response.WithdrawResFailure(proto.T_FailureCode_T_GENERAL_ERROR, "Something went wrong"), nil
}

func (s DefaultBankingService) Deposit(ctx context.Context, req *proto.DepositReq) (res *proto.DepositRes, err error) {
	if req.Amount == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "amount is empty")
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return response.DepositResFailure(proto.T_FailureCode_T_MISSING_DATA, "userId is empty"), nil
		}

		kMessage := bson.D{
			{Key: "debtor", Value: nil},
			{Key: "amount", Value: req.Amount},
			{Key: "creditor", Value: userId[0]},
			{Key: "event", Value: "Deposit"},
		}

		kMessageBytes, err := bson.Marshal(kMessage)
		if err != nil {
			return response.DepositResFailure(proto.T_FailureCode_T_GENERAL_ERROR, "Something went wrong"), nil
		}

		err = s.kafkaProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &s.kafkaTopic,
				Partition: kafka.PartitionAny,
			},
			Value: kMessageBytes,
		}, nil)

		if err != nil {
			return response.DepositResFailure(proto.T_FailureCode_T_NETWORK_ERROR, "Kafka network error"), nil
		}

		return response.DepositResSuccess(fmt.Sprintf("You deposited $%.2f in your balance", req.Amount)), nil
	}

	return response.DepositResFailure(proto.T_FailureCode_T_GENERAL_ERROR, "Something went wrong"), nil
}

func (s DefaultBankingService) Transfer(ctx context.Context, req *proto.TransferReq) (res *proto.TransferRes, err error) {
	if req.Amount == 0 || req.To == "" {
		return nil, status.Errorf(codes.InvalidArgument, "incomplete request")
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return response.TransferResFailure(proto.T_FailureCode_T_MISSING_DATA, "userId is empty"), nil
		}

		md := metadata.New(map[string]string{"userId": userId[0]})
		ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

		existingBalance, err := s.grpcClient.Client.BalanceInquiry(ctxWithMd, &paymentProto.BalanceReq{})

		if err != nil {
			return response.TransferResFailure(proto.T_FailureCode_T_NETWORK_ERROR, "Error calling the rpc"), nil
		}

		switch result := existingBalance.Result.(type) {
		case *paymentProto.BalanceRes_Success_:
			if result.Success.Balance < req.Amount {
				return response.TransferResFailure(proto.T_FailureCode_T_INSUFFICIENT_BALANCE, "Insufficient funds"), nil
			}

			kMessage := bson.D{
				{Key: "debtor", Value: userId[0]},
				{Key: "amount", Value: req.Amount},
				{Key: "creditor", Value: req.To},
				{Key: "event", Value: "P2PTransfer"},
			}

			kMessageBytes, err := bson.Marshal(kMessage)
			if err != nil {
				return response.TransferResFailure(proto.T_FailureCode_T_GENERAL_ERROR, "Something went wrong"), nil
			}

			err = s.kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &s.kafkaTopic,
					Partition: kafka.PartitionAny,
				},
				Value: kMessageBytes,
			}, nil)
			return response.TransferResSuccess(fmt.Sprintf("You have transferred $%.2f successfully.", req.Amount)), nil
		case *paymentProto.BalanceRes_Failure:
			if result.Failure.FailureCode == paymentProto.FailureCode_MISSING_DATA {
				return response.TransferResFailure(proto.T_FailureCode_T_MISSING_DATA, "Missing data"), nil

			} else if result.Failure.FailureCode == paymentProto.FailureCode_INSUFFICIENT_BALANCE {
				return response.TransferResFailure(proto.T_FailureCode_T_INSUFFICIENT_BALANCE, "Insufficient Balance"), nil
			}

			return response.TransferResFailure(proto.T_FailureCode_T_GENERAL_ERROR, "Something went wrong"), nil
		}
	}

	return response.TransferResFailure(proto.T_FailureCode_T_GENERAL_ERROR, "Something went wrong"), nil
}
