package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
)

type entry struct {
	Data string `json:"data" binding:"required"`
}

type response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

var collectionName = "halan"

func addTodo(c *gin.Context) {
	var newTodo entry

	collc := Collc(collectionName)

	if err := c.BindJSON(&newTodo); err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(http.StatusCreated, res)
		return
	}

	isInserted := InsertTodo(collc, newTodo)

	if isInserted {
		res := &response{
			Message: "success",
			Status:  http.StatusOK,
			Data:    newTodo,
		}

		c.IndentedJSON(res.Status, res)
	} else {

		res := &response{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    nil,
		}

		c.IndentedJSON(res.Status, res)

	}
}

func getTodos(c *gin.Context) {
	collc := Collc(collectionName)
	data := FilterTodos(collc, nil)

	res := &response{
		Message: "success",
		Status:  http.StatusOK,
		Data:    data,
	}

	c.IndentedJSON(res.Status, res)
}

func getTodoById(c *gin.Context) {
	id := c.Param("id")

	collc := Collc(collectionName)

	objId, _ := primitive.ObjectIDFromHex(id)

	data := FilterTodos(collc, &bson.D{{
		"_id", objId,
	}})

	if len(data) > 0 {
		res := &response{
			Message: "success",
			Status:  http.StatusOK,
			Data:    data[0],
		}

		c.IndentedJSON(http.StatusCreated, res)
		return

	}

	res := &response{
		Message: "fail",
		Status:  http.StatusNotFound,
		Data:    "Item not found!",
	}

	c.IndentedJSON(res.Status, res)

}

func removeTodoById(c *gin.Context) {
	id := c.Param("id")

	i, err := primitive.ObjectIDFromHex(id)

	collc := Collc(collectionName)

	if err != nil {
		res := &response{
			Message: "fail",
			Status:  http.StatusInternalServerError,
			Data:    "Something went wrong",
		}

		c.IndentedJSON(http.StatusCreated, res)
		return
	}

	isDeleted := DeleteTodo(collc, i)

	if isDeleted {
		res := &response{
			Message: "success",
			Status:  http.StatusOK,
			Data:    nil,
		}

		c.IndentedJSON(http.StatusCreated, res)

		return
	}

	res := &response{
		Message: "fail",
		Status:  http.StatusInternalServerError,
		Data:    "Item not found!",
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func FilterTodos(collection *mongo.Collection, query *bson.D) []bson.M {

	if query == nil {
		query = &bson.D{}
	}

	cursor, err := collection.Find(context.Background(), query)

	if err != nil {
		fmt.Println("Something went wrong: ", err)
		return nil
	}

	var results []bson.M

	if err = cursor.All(context.Background(), &results); err != nil {
		fmt.Println("Something went wrong: ", err)
		return nil
	}

	return results
}

func InsertTodo(collection *mongo.Collection, data entry) bool {

	cursor, err := collection.InsertOne(context.Background(), data)

	if err != nil {
		fmt.Println("Something went wrong: ", err)
		return false
	}

	fmt.Println("Inserted: ", cursor.InsertedID)

	return true

}

func DeleteTodo(collection *mongo.Collection, entryId primitive.ObjectID) bool {
	cursor, err := collection.DeleteOne(context.TODO(), bson.D{
		{"_id", entryId},
	})

	if err != nil {
		fmt.Println("Something went wrong: ", err)
		return false
	}

	if cursor.DeletedCount > 0 {
		return true
	}

	return false

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
