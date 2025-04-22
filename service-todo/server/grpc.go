package server

import (
	"assm/service-todo/ctx"
	"assm/service-todo/proto"
	"assm/service-todo/service"
	"context"
)

type DefaultTodoGRPCServer struct {
	proto.UnimplementedTodoServiceServer
	serviceContext ctx.ServiceCtx
	todoService    *service.DefaultBankingService
}

func NewDefaultTodoGRPCServer(serviceContext ctx.ServiceCtx) *DefaultTodoGRPCServer {
	return &DefaultTodoGRPCServer{
		serviceContext: serviceContext,
		todoService:    service.NewDefaultBankingService(serviceContext),
	}
}

func (s *DefaultTodoGRPCServer) Withdraw(ctx context.Context, req *proto.WithdrawReq) (res *proto.WithdrawRes, err error) {
	return s.todoService.Withdraw(ctx, req)
}

func (s DefaultTodoGRPCServer) Deposit(ctx context.Context, req *proto.DepositReq) (res *proto.DepositRes, err error) {
	return s.todoService.Deposit(ctx, req)
}

func (s DefaultTodoGRPCServer) Transfer(ctx context.Context, req *proto.TransferReq) (res *proto.TransferRes, err error) {
	return s.todoService.Transfer(ctx, req)
}
