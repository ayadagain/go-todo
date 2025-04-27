package response

import "assm/service-todo/proto"

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
