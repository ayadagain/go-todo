package service

import (
	"assm/service-payment/ctx"
	"assm/service-payment/proto"
	"assm/service-todo/response"
	"context"
	"google.golang.org/grpc/metadata"
)

type DefaultPaymentService struct {
	ServiceContext ctx.ServiceCtx
}

func NewDefaultService(serviceCtx ctx.ServiceCtx) *DefaultPaymentService {
	return &DefaultPaymentService{
		ServiceContext: serviceCtx,
	}
}

func (s DefaultPaymentService) BalanceInquiry(ctx context.Context, _ *proto.BalanceReq) (res *proto.BalanceRes, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return response.BalanceResFailure(proto.FailureCode_MISSING_DATA, "User Id is empty"), nil
		}

		existingBalance, err := s.GetBalance(userId[0])

		if err != nil {
			return response.BalanceResFailure(proto.FailureCode_GENERAL_ERROR, err.Error()), nil
		}

		return response.BalanceResSuccess(float32(existingBalance)), nil
	}

	return response.BalanceResFailure(proto.FailureCode_MISSING_DATA, "User Id is empty"), nil
}
