package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ConnectDB(uri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		panic("Failed to create client: " + err.Error())
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

type entry struct {
	ID   int    `json:"id" binding:"required"`
	Data string `json:"data" binding:"required"`
}

type response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

var todos = []entry{}

func addTodo(c *gin.Context) {
	var newTodo entry

	if err := c.BindJSON(&newTodo); err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(http.StatusCreated, res)
		return
	}

	todos = append(todos, newTodo)

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data:    newTodo,
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func getTodos(c *gin.Context) {

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data:    todos,
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func getTodoById(c *gin.Context) {
	id := c.Param("id")

	i, err := strconv.Atoi(id)

	if err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusInternalServerError,
			Data:    "Something went wrong",
		}

		c.IndentedJSON(http.StatusCreated, res)
		return
	}

	for _, item := range todos {
		if item.ID == i {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}

	res := &response{
		Message: "fail",
		Status:  http.StatusNotFound,
		Data:    "Item not found!",
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func removeTodoById(c *gin.Context) {
	id := c.Param("id")

	i, err := strconv.Atoi(id)

	if err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusInternalServerError,
			Data:    "Something went wrong",
		}

		c.IndentedJSON(http.StatusCreated, res)
		return
	}

	for key, item := range todos {
		if item.ID == i {
			todos = append(todos[:key], todos[key+1:]...)

			res := &response{
				Message: "fail",
				Status:  http.StatusInternalServerError,
				Data:    "Item deleted successfully",
			}

			c.IndentedJSON(http.StatusCreated, res)

			return
		}
	}

	res := &response{
		Message: "fail",
		Status:  http.StatusInternalServerError,
		Data:    "Item not found!",
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collec := client.Database("halan").Collection(collectionName)
	return collec
}

func main() {

	_ = godotenv.Load()

	mongoUri := os.Getenv("MONGO_URI")
	collectionName := os.Getenv("MONGO_COLLECTION_NAME")

	if mongoUri == "" {
		panic("env variable MONGO_URI not found")
	}

	if collectionName == "" {
		panic("env variable MONGO_COLLECTION_NAME not found")
	}

	var DB *mongo.Client = ConnectDB(mongoUri)

	cursor, err := GetCollection(DB, collectionName).Find(context.Background(), bson.D{{}})

	if err != nil {
		fmt.Println(err)
		return
	}

	var collectionData []bson.M

	if err = cursor.All(context.Background(), &collectionData); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("collectionData:", collectionData)

	router := gin.Default()

	router.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome to the api",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodoById)
	router.DELETE("/todos/:id", removeTodoById)
	router.POST("/todos", addTodo)

	router.Run() // listen and serve on 0.0.0.0:8080
}
