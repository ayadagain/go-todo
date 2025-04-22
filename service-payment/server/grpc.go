package server

import (
	"assm/service-payment/ctx"
	"assm/service-payment/proto"
	"assm/service-payment/service"
	"context"
	"fmt"
)

type DefaultPaymentGRPCServer struct {
	proto.UnimplementedPaymentServiceServer
	serviceContext ctx.ServiceCtx
	paymentService *service.DefaultPaymentService
}

func NewDefaultPaymentGRPCServer(serviceContext ctx.ServiceCtx) *DefaultPaymentGRPCServer {
	return &DefaultPaymentGRPCServer{
		serviceContext: serviceContext,
		paymentService: service.NewDefaultService(serviceContext),
	}
}

func (o *DefaultPaymentGRPCServer) BalanceInquiry(ctx context.Context, req *proto.BalanceReq) (res *proto.BalanceRes, err error) {
	fmt.Println("Received request: ", req)
	return o.paymentService.BalanceInquiry(ctx, req)
}
