package grpc

import (
	"assm/service-payment/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GrpcConn() proto.PaymentServiceClient {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	c := proto.NewPaymentServiceClient(conn)
	return c
}
