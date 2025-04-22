package service

import (
	"assm/service-payment/ctx"
	"assm/service-payment/proto"
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

func (s DefaultPaymentService) BalanceInquiry(ctx context.Context, req *proto.BalanceReq) (res *proto.BalanceRes, err error) {
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

		existingBalance, err := s.GetBalance(userId[0])

		if err != nil {
			return &proto.BalanceRes{
				Result: &proto.BalanceRes_Failure{
					Failure: &proto.Failure{
						FailureCode:    proto.FailureCode_GENERAL_ERROR,
						FailureMessage: err.Error(),
					},
				},
			}, nil
		}

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
