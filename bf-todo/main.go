package main

import (
	hgrpc "assm/bf-todo/grpc"
	"assm/service-todo/proto"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"os"
)

type response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

type rpcResponse struct {
	ObjectId  string `json:"_id,omitempty" bson:"_id"`
	Message   string `json:"message"`
	CreatedBy string `json:"created_by,omitempty"`
}

type networking struct {
	Hgrpc proto.TodoServiceClient
}

type transaction struct {
	Amount float32 `json:"amount"`
	To     string  `json:"to"`
}

func newNetworking() *networking {
	return &networking{
		Hgrpc: hgrpc.GrpcConn(),
	}
}

func auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		serverKey := c.Request.Header.Get("ServerKey")

		if serverKey != secretKey {
			res := &response{
				Message: "fail",
				Status:  http.StatusUnauthorized,
				Data:    nil,
			}

			c.IndentedJSON(res.Status, res)
			c.Abort()

			return
		}

		c.Next()
	}
}

func (s *networking) withdraw(c *gin.Context) {
	var postData transaction

	if err := c.BindJSON(&postData); err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	if postData.Amount < 0 {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    "Cannot send negative amount",
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	r, err := s.Hgrpc.Withdraw(ctxWithMd, &proto.WithdrawReq{
		Amount: postData.Amount,
	})

	if err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data:    r,
	}

	c.IndentedJSON(res.Status, res)

}

func (s *networking) deposit(c *gin.Context) {
	var postData transaction

	if err := c.BindJSON(&postData); err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	if postData.Amount < 0 {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    "Cannot send negative amount",
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	r, err := s.Hgrpc.Deposit(ctxWithMd, &proto.DepositReq{
		Amount: postData.Amount,
	})

	if err != nil {
		log.Println("err: ", err)
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data:    r,
	}

	c.IndentedJSON(res.Status, res)

}

func (s *networking) transfer(c *gin.Context) {
	var postData transaction

	if err := c.BindJSON(&postData); err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	if postData.To == "" {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    "Missing To field",
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	md := metadata.New(map[string]string{"userId": "67fbd5d9fc7128b743d265b7"})
	ctxWithMd := metadata.NewOutgoingContext(context.Background(), md)

	r, err := s.Hgrpc.Transfer(ctxWithMd, &proto.TransferReq{
		Amount: postData.Amount,
		To:     postData.To,
	})

	if err != nil {
		log.Println("err: ", err)
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data:    r,
	}

	c.IndentedJSON(res.Status, res)

}

func main() {
	_ = godotenv.Load("../.env")

	secretKey := os.Getenv("SECRET_KEY")

	router := gin.Default()
	net := newNetworking()

	router.Use(auth(secretKey))

	router.POST("/withdraw", net.withdraw)
	router.POST("/deposit", net.deposit)
	router.POST("/transfer", net.transfer)

	err := router.Run()

	if err != nil {
		panic("failed to run the server")
	}
}
