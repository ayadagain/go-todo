package main

import (
	hgrpc "assm/bf-todo/grpc"
	"assm/service-todo/proto"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type entry struct {
	Data string `json:"data" binding:"required" bson:"data"`
}

type response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

type rpcResponse struct {
	ObjectId string `json:"_id,omitempty" bson:"_id"`
	Message  string `json:"message"`
}

type networking struct {
	Hgrpc proto.TodoServiceClient
}

func newNetworking() *networking {
	return &networking{
		Hgrpc: hgrpc.GrpcConn(),
	}
}

func (s *networking) getTodos(c *gin.Context) {
	ctx := context.TODO()
	var result []rpcResponse

	r, err := s.Hgrpc.GetTodos(ctx, &proto.GetTodosReq{})

	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}

	for _, todo := range r.Todos {
		result = append(result, rpcResponse{
			ObjectId: todo.ObjectId,
			Message:  todo.Message,
		})
	}

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data:    result,
	}

	c.IndentedJSON(res.Status, res)
}

func (s *networking) getTodoById(c *gin.Context) {
	ctx := context.TODO()

	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		res := &response{
			Message: "Invalid id",
			Status:  http.StatusBadRequest,
			Data:    nil,
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	objIdString := objId.Hex()

	r, err := s.Hgrpc.GetTodo(ctx, &proto.GetTodoReq{
		ObjectId: objIdString,
	})

	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data: &rpcResponse{
			ObjectId: r.ObjectId,
			Message:  r.Message,
		},
	}

	c.IndentedJSON(res.Status, res)
}

func (s *networking) removeTodoById(c *gin.Context) {
	ctx := context.TODO()

	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		res := &response{
			Message: "Invalid id",
			Status:  http.StatusBadRequest,
			Data:    nil,
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	objIdString := objId.Hex()

	_, err = s.Hgrpc.DeleteTodo(ctx, &proto.DeleteTodoReq{
		ObjectId: objIdString,
	})

	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data:    nil,
	}

	c.IndentedJSON(res.Status, res)
	return
}

func (s *networking) insertTodo(c *gin.Context) {
	ctx := context.TODO()

	var newTodo entry

	if err := c.BindJSON(&newTodo); err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	r, err := s.Hgrpc.InsertTodo(ctx, &proto.InsertTodoReq{
		Data: newTodo.Data,
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
		Data:    &rpcResponse{Message: r.Data},
	}

	c.IndentedJSON(res.Status, res)
}

func main() {
	router := gin.Default()
	net := newNetworking()

	router.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome to the apiÏ",
		})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/todos", net.getTodos)
	router.GET("/todos/:id", net.getTodoById)
	router.POST("/todos", net.insertTodo)
	router.DELETE("/todos/:id", net.removeTodoById)

	err := router.Run()

	if err != nil {
		panic("failed to run the server")
	}
}
