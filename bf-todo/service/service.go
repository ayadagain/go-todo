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

var (
	ErrFailedTransaction = errors.New("fail to establish transaction")
	ErrSmthWentWrong     = errors.New("something went wrong")
	ErrNegVals           = errors.New("cannot send negative amount")
	ErrInSufficientFunds = errors.New("insufficient funds")
	ErrMissingData       = errors.New("missing data")
)

type DefaultService struct {
	ServiceContext ctx.ServiceCtx
	TodoGrpcClient *client.TodoClient
}

func NewDefaultService(serviceContext ctx.ServiceCtx) *DefaultService {
	return &DefaultService{
		ServiceContext: serviceContext,
		TodoGrpcClient: client.NewTodoClient(serviceContext.GrpcClient()),
	}
}

func (a *DefaultService) Withdraw(amount float32) (response *model.TransactionResponse, err error) {
	if amount < 0 {
		return nil, ErrFailedTransaction
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	res, err := a.TodoGrpcClient.Client.Withdraw(ctxWithMd, &proto.WithdrawReq{Amount: amount})
	if err != nil {
		return nil, ErrSmthWentWrong
	}

	switch result := res.Result.(type) {
	case *proto.WithdrawRes_Success_:
		return &model.TransactionResponse{
			Status:  int(result.Success.Status),
			Message: result.Success.Message,
		}, nil

	case *proto.WithdrawRes_Failure:
		if result.Failure.FailureCode == proto.T_FailureCode_T_MISSING_DATA {
			return nil, ErrMissingData
		} else if result.Failure.FailureCode == proto.T_FailureCode_T_INSUFFICIENT_BALANCE {
			return nil, ErrInSufficientFunds
		}

		return nil, ErrSmthWentWrong
	}

	return nil, ErrSmthWentWrong
}

func (a *DefaultService) Deposit(amount float32) (response *model.TransactionResponse, err error) {
	if amount < 0 {
		return nil, ErrNegVals
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	res, err := a.TodoGrpcClient.Client.Deposit(ctxWithMd, &proto.DepositReq{Amount: amount})

	if err != nil {
		return nil, ErrSmthWentWrong
	}

	switch result := res.Result.(type) {
	case *proto.DepositRes_Success_:
		return &model.TransactionResponse{
			Status:  int(result.Success.Status),
			Message: result.Success.Message,
		}, nil
	case *proto.DepositRes_Failure:
		return nil, ErrSmthWentWrong
	default:
		return nil, ErrSmthWentWrong
	}
}

func (a *DefaultService) Transfer(to string, amount float32) (response *model.TransactionResponse, err error) {
	if to == "" {
		return nil, ErrMissingData
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	res, err := a.TodoGrpcClient.Client.Transfer(ctxWithMd, &proto.TransferReq{
		Amount: amount,
		To:     to,
	})

	if err != nil {
		return nil, ErrSmthWentWrong
	}

	switch result := res.Result.(type) {
	case *proto.TransferRes_Success_:
		return &model.TransactionResponse{
			Status:  int(result.Success.Status),
			Message: result.Success.Message,
		}, nil
	case *proto.TransferRes_Failure:
		if result.Failure.FailureCode == proto.T_FailureCode_T_MISSING_DATA {
			return nil, ErrMissingData
		} else if result.Failure.FailureCode == proto.T_FailureCode_T_INSUFFICIENT_BALANCE {
			return nil, ErrInSufficientFunds
		}
		return nil, ErrSmthWentWrong
	default:
		return nil, ErrSmthWentWrong
	}
}
