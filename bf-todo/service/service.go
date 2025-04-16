package service

import (
	"assm/bf-todo/ctx"
	"assm/bf-todo/grpc/client"
	"assm/bf-todo/model"
	"assm/service-todo/proto"
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

type DefaultService struct {
	ServiceContext ctx.ServiceCtx
	TodoGrpcClient *client.TodoClient
}

func NewDefaultService(serviceContext ctx.ServiceCtx) *DefaultService {
	todoGrpcClient := client.NewTodoClient(serviceContext.GrpcClient())
	return &DefaultService{
		ServiceContext: serviceContext,
		TodoGrpcClient: todoGrpcClient,
	}
}

func (a *DefaultService) Withdraw(amount float32) (response *model.TransactionResponse, err error) {
	if amount < 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	res, err := a.TodoGrpcClient.Client.Withdraw(ctxWithMd, &proto.WithdrawReq{Amount: amount})
	if err != nil {
		return nil, errors.New("fail to establish transaction")
	}

	switch result := res.Result.(type) {
	case *proto.WithdrawRes_Success_:
		return &model.TransactionResponse{
			Status:  int(result.Success.Status),
			Message: result.Success.Message,
		}, nil

	case *proto.WithdrawRes_Failure:
		return nil, errors.New(result.Failure.FailureMessage)
	}

	return nil, errors.New("Something went wrong")

}

func (a *DefaultService) Deposit(amount float32) (response *model.TransactionResponse, err error) {
	if amount < 0 {
		return &model.TransactionResponse{
			Status:  0,
			Message: "Cannot send negative amount",
		}, nil
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	res, err := a.TodoGrpcClient.Client.Deposit(ctxWithMd, &proto.DepositReq{Amount: amount})

	if err != nil {
		return &model.TransactionResponse{
			Status:  0,
			Message: "Something went wrong",
		}, nil
	}

	switch result := res.Result.(type) {
	case *proto.DepositRes_Success_:
		return &model.TransactionResponse{
			Status:  int(result.Success.Status),
			Message: result.Success.Message,
		}, nil
	case *proto.DepositRes_Failure:
		return &model.TransactionResponse{
			Status:  int(result.Failure.FailureCode),
			Message: result.Failure.FailureMessage,
		}, nil

	default:
		return nil, errors.New("Something went wrong")

	}
}

func (a *DefaultService) Transfer(to string, amount float32) (response *model.TransactionResponse, err error) {

	if to == "" {
		return nil, errors.New("To field is required")
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	res, err := a.TodoGrpcClient.Client.Transfer(ctxWithMd, &proto.TransferReq{
		Amount: amount,
		To:     to,
	})

	if err != nil {
		return &model.TransactionResponse{
			Status:  -1,
			Message: "Something went wrong",
		}, nil
	}

	switch result := res.Result.(type) {
	case *proto.TransferRes_Success_:
		return &model.TransactionResponse{
			Status:  int(result.Success.Status),
			Message: result.Success.Message,
		}, nil
	case *proto.TransferRes_Failure:
		return nil, errors.New(result.Failure.FailureMessage)

	default:
		return nil, errors.New("Something went wrong")

	}
}
