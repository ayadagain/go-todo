package main

import (
	"assm/service-todo/ctx"
	"assm/service-todo/kafka/producer"
	"assm/service-todo/proto"
	"assm/service-todo/server"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func createGRPCServer(serviceContext ctx.ServiceCtx) *grpc.Server {
	grpcServer := grpc.NewServer()
	grpcService := server.NewDefaultTodoGRPCServer(serviceContext)
	proto.RegisterTodoServiceServer(grpcServer, grpcService)
	return grpcServer
}

func main() {

	_ = godotenv.Load("../.env")
	serviceCtx := ctx.NewDefaultServiceCtx()
	grpcServer := createGRPCServer(serviceCtx)
	kafkaProducer := producer.NewTodoProducer(serviceCtx)

	serviceCtx.Start(grpcServer, kafkaProducer)
	serviceCtx.ShutdownHook(kafkaProducer.Shutdown())
	fmt.Println("Shutting down...")
}
