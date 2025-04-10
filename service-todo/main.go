package main

import (
	"assm/service-todo/db"
	"assm/service-todo/proto"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type Server struct {
	proto.UnimplementedTodoServiceServer
	dbCollection *mongo.Collection
}

func (s *Server) GetTodo(ctx context.Context, req *proto.GetTodoReq) (res *proto.GetTodoRes, err error) {
	var qry *bson.D

	if req.ObjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "ObjectId is required")
	}

	oId, err := primitive.ObjectIDFromHex(req.ObjectId)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	qry = &bson.D{{
		"_id", oId,
	}}

	data := db.FilterTodos(s.dbCollection, qry)

	if len(data) == 0 || data == nil {
		return nil, status.Errorf(codes.NotFound, "not found")
	}

	resp := data[0]

	protoReq := &proto.GetTodoRes{
		ObjectId:  resp.ID.Hex(),
		Message:   resp.Data,
		CreatedBy: resp.CreatedBy,
	}
	return protoReq, nil
}

func (s *Server) GetTodos(ctx context.Context, req *proto.GetTodosReq) (res *proto.GetTodosRes, err error) {
	var pp []*proto.GetTodoRes
	var qry *bson.D

	data := db.FilterTodos(s.dbCollection, qry)

	if len(data) == 0 || data == nil {
		return nil, status.Errorf(codes.NotFound, "not found")
	}

	for _, t := range data {
		_id := t.ID.Hex()

		pp = append(pp, &proto.GetTodoRes{
			Message:   t.Data,
			ObjectId:  _id,
			CreatedBy: t.CreatedBy,
		})
	}

	protoReq := &proto.GetTodosRes{Todos: pp}
	return protoReq, nil
}

func (s *Server) DeleteTodo(ctx context.Context, req *proto.DeleteTodoReq) (res *proto.DeleteTodoRes, err error) {
	if req.ObjectId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "objectId is required")
	}

	if req.ObjectId != "" {
		oId, _ := primitive.ObjectIDFromHex(req.ObjectId)

		isDeleted := db.DeleteTodo(s.dbCollection, oId)

		if isDeleted {
			return &proto.DeleteTodoRes{}, nil
		}

		return nil, status.Errorf(codes.NotFound, "not found")

	}

	return nil, status.Errorf(codes.NotFound, "Something went wrong")
}

func (s *Server) InsertTodo(ctx context.Context, req *proto.InsertTodoReq) (res *proto.InsertTodoRes, err error) {
	if req.Data == "" {
		return nil, status.Errorf(codes.InvalidArgument, "body is empty")
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "userId is empty")
		}

		isInserted := db.InsertTodo(s.dbCollection, req.Data, userId[0])

		if isInserted {
			return &proto.InsertTodoRes{
				Data: req.Data,
			}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "Something went wrong")
}

func main() {
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()

	proto.RegisterTodoServiceServer(s, &Server{
		dbCollection: db.Collc("halan"),
	})

	log.Printf("server listening at %v", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
