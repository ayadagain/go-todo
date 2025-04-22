package client

import (
	"assm/service-payment/proto"
	"google.golang.org/grpc"
)

type PaymentClient struct {
	Client proto.PaymentServiceClient
}

func NewPaymentClient(grpcConnection *grpc.ClientConn) *PaymentClient {
	return &PaymentClient{Client: proto.NewPaymentServiceClient(grpcConnection)}
}
