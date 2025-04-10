package grpc

import (
	"assm/service-todo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GrpcConn() proto.TodoServiceClient {
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	c := proto.NewTodoServiceClient(conn)
	return c
}
