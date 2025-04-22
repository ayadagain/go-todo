package main

import (
	"assm/service-payment/ctx"
	"assm/service-payment/kafka/consumer"
	"assm/service-payment/proto"
	"assm/service-payment/server"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func createGRPCServer(serviceContext ctx.ServiceCtx) *grpc.Server {
	grpcServer := grpc.NewServer()
	grpcService := server.NewDefaultPaymentGRPCServer(serviceContext)
	proto.RegisterPaymentServiceServer(grpcServer, grpcService)
	return grpcServer
}

func main() {
	_ = godotenv.Load("../.env")
	serviceContext := ctx.NewDefaultServiceCtx()
	grpcServer := createGRPCServer(serviceContext)
	kafkaConsumer := consumer.NewPaymentConsumer(serviceContext)
	serviceContext.Start(grpcServer, kafkaConsumer)

	serviceContext.ShutdownHook()
	fmt.Println("Shutting down...")
}
