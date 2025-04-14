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
	client            *mongo.Client
	dbCollection      *mongo.Collection
	paymentCollection *mongo.Collection
	userCollection    *mongo.Collection
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

		_, err := db.Transfer(s.client, s.userCollection, "67fbd5befc7128b743d265b6", "67fbd5d9fc7128b743d265b7", 100)

		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

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

func (s *Server) Withdraw(ctx context.Context, req *proto.WithdrawReq) (res *proto.WithdrawRes, err error) {
	if req.Amount == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "amount is empty")
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "userId is empty")
		}

		existingBalance := db.GetBalance(s.userCollection, userId[0])

		if existingBalance < float64(req.Amount) {
			return nil, status.Errorf(codes.NotFound, "Balance is insufficient")
		}

		insertOp := db.EditBalance(s.client, s.userCollection, userId[0], req.Amount*-1)

		if insertOp {
			return &proto.WithdrawRes{
				Status:  1,
				Message: fmt.Sprintf("You withdrew $%.2f from your balance. Your balance now is $%.2f", req.Amount, float32(existingBalance)-req.Amount),
			}, nil
		}

		return nil, status.Errorf(codes.NotFound, "Something went wrong")

	}

	return nil, status.Errorf(codes.NotFound, "Something went wrong")
}

func (s *Server) Deposit(ctx context.Context, req *proto.DepositReq) (res *proto.DepositRes, err error) {
	if req.Amount == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "amount is empty")
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "userId is empty")
		}

		insertOp := db.EditBalance(s.client, s.userCollection, userId[0], req.Amount)

		if insertOp {
			return &proto.DepositRes{
				Status:  1,
				Message: fmt.Sprintf("You deposited $%.2f in your balance", req.Amount),
			}, nil
		}

		return nil, status.Errorf(codes.NotFound, "Something went wrong")

	}

	return nil, status.Errorf(codes.NotFound, "Something went wrong")
}

func (s *Server) Transfer(ctx context.Context, req *proto.TransferReq) (res *proto.TransferRes, err error) {
	if req.Amount == 0 || req.To == "" {
		return nil, status.Errorf(codes.InvalidArgument, "incomplete request")
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userId := md["userid"]

		if len(userId) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "userId is empty")
		}

		insertOp, err := db.Transfer(s.client, s.userCollection, userId[0], req.To, req.Amount)

		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		if insertOp {
			return &proto.TransferRes{
				Status:  1,
				Message: fmt.Sprintf("You have transferred $%.2f successfully.", req.Amount),
			}, nil
		}

		return nil, status.Errorf(codes.NotFound, "Something went wrong")
	}

	return nil, status.Errorf(codes.NotFound, "Something went wrong")

	return &proto.TransferRes{}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()

	mongoClient := db.Conn()

	proto.RegisterTodoServiceServer(s, &Server{
		client:            mongoClient,
		dbCollection:      db.Collc(mongoClient, "halan"),
		paymentCollection: db.Collc(mongoClient, "payments"),
		userCollection:    db.Collc(mongoClient, "users"),
	})

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
