package response

import (
	"assm/bf-todo/model"
	paymentProto "assm/service-payment/proto"
	"assm/service-todo/proto"
)

func WithdrawResSuccess(message string) *proto.WithdrawRes {
	return &proto.WithdrawRes{
		Result: &proto.WithdrawRes_Success_{Success: &proto.WithdrawRes_Success{
			Message: message,
		}},
	}
}

func WithdrawResFailure(failureCode proto.T_FailureCode, failureMessage string) *proto.WithdrawRes {
	return &proto.WithdrawRes{
		Result: &proto.WithdrawRes_Failure{Failure: &proto.T_Failure{
			FailureCode:    failureCode,
			FailureMessage: failureMessage,
		}},
	}
}

func DepositResSuccess(message string) *proto.DepositRes {
	return &proto.DepositRes{
		Result: &proto.DepositRes_Success_{Success: &proto.DepositRes_Success{
			Message: message,
		}},
	}
}

func DepositResFailure(failureCode proto.T_FailureCode, failureMessage string) *proto.DepositRes {
	return &proto.DepositRes{
		Result: &proto.DepositRes_Failure{Failure: &proto.T_Failure{
			FailureCode:    failureCode,
			FailureMessage: failureMessage,
		}},
	}
}

func TransferResSuccess(message string) *proto.TransferRes {
	return &proto.TransferRes{
		Result: &proto.TransferRes_Success_{Success: &proto.TransferRes_Success{
			Message: message,
		}},
	}
}

func TransferResFailure(failureCode proto.T_FailureCode, failureMessage string) *proto.TransferRes {
	return &proto.TransferRes{
		Result: &proto.TransferRes_Failure{Failure: &proto.T_Failure{
			FailureCode:    failureCode,
			FailureMessage: failureMessage,
		}},
	}
}

func TransactionResponse(message string) *model.TransactionResponse {
	return &model.TransactionResponse{
		Message: message,
	}
}

func BalanceResSuccess(balance float32) *paymentProto.BalanceRes {
	return &paymentProto.BalanceRes{
		Result: &paymentProto.BalanceRes_Success_{
			Success: &paymentProto.BalanceRes_Success{
				Balance: balance,
			},
		},
	}
}

func BalanceResFailure(failureCode paymentProto.FailureCode, failureMessage string) *paymentProto.BalanceRes {
	return &paymentProto.BalanceRes{
		Result: &paymentProto.BalanceRes_Failure{
			Failure: &paymentProto.Failure{
				FailureCode:    failureCode,
				FailureMessage: failureMessage,
			},
		},
	}
}
