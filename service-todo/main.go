package main

import (
	"assm/service-todo/db"
	"assm/service-todo/proto"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

var collectionName = "halan"

type Server struct {
	proto.UnimplementedTodoServiceServer
}

func (s *Server) SelectTodo(ctx context.Context, req *proto.TodoReq) (res *proto.Todos, err error) {
	var pp []*proto.TodoRes
	var qry *bson.D

	collc := db.Collc(collectionName)

	if req.ObjectId != nil {
		oId, _ := primitive.ObjectIDFromHex(*req.ObjectId)

		qry = &bson.D{{
			"_id", oId,
		}}
	}

	data := db.FilterTodos(collc, qry)

	if len(data) == 0 || data == nil {
		return nil, status.Errorf(codes.NotFound, "not found")
	}

	for _, t := range data {
		_id := t["_id"].(primitive.ObjectID).Hex()

		pp = append(pp, &proto.TodoRes{
			Message:  t["data"].(string),
			ObjectId: &_id,
		})
	}

	protoReq := &proto.Todos{Todos: pp}
	return protoReq, nil
}

func (s *Server) DeleteTodo(ctx context.Context, req *proto.TodoReq) (res *proto.TodoRes, err error) {
	collc := db.Collc(collectionName)

	if *req.ObjectId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "objectId is required")
	}

	if req.ObjectId != nil {
		oId, _ := primitive.ObjectIDFromHex(*req.ObjectId)

		isDeleted := db.DeleteTodo(collc, oId)

		if isDeleted {
			return &proto.TodoRes{
				Message: "",
			}, nil
		}

		return nil, status.Errorf(codes.NotFound, "not found")

	}

	return nil, status.Errorf(codes.NotFound, "Something went wrong")
}

func (s *Server) InsertTodo(ctx context.Context, req *proto.TodoReq) (res *proto.TodoRes, err error) {
	collc := db.Collc(collectionName)

	if req.Body == "" {
		return nil, status.Errorf(codes.InvalidArgument, "body is empty")
	}

	isInserted := db.InsertTodo(collc, req.Body)

	if isInserted {
		return &proto.TodoRes{
			Message: req.Body,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "Something went wrong")
}
func main() {
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()

	proto.RegisterTodoServiceServer(s, &Server{})

	log.Printf("server listening at %v", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
