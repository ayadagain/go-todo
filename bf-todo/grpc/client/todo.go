package client

import (
	"assm/service-todo/proto"
	"google.golang.org/grpc"
)

type TodoClient struct {
	Client proto.TodoServiceClient
}

func NewTodoClient(grpcConnection *grpc.ClientConn) *TodoClient {
	return &TodoClient{Client: proto.NewTodoServiceClient(grpcConnection)}
}
