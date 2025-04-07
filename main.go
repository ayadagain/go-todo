package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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

func main() {
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
