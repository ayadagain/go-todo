package route

import (
	"assm/bf-todo/ctx"
	"assm/bf-todo/model"
	"assm/bf-todo/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
	serviceContext ctx.ServiceCtx
	todoService    *service.DefaultService
}

func NewRouter(serviceCtx ctx.ServiceCtx) *Router {
	ser := service.NewDefaultService(serviceCtx)
	return &Router{
		serviceContext: serviceCtx,
		todoService:    ser,
	}
}

func (r *Router) Install(ginEngine *gin.Engine) {
	router := ginEngine.Group("")
	router.Use(r.ServerKeyCheck)
	router.POST("/withdraw", r.withdraw)
	router.POST("/deposit", r.deposit)
	router.POST("/transfer", r.transfer)
}

func (r *Router) withdraw(c *gin.Context) {
	var postData model.Transaction

	if err := c.BindJSON(&postData); err != nil {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(http.StatusOK, res)
		return
	}

	if postData.Amount < 0 {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    "Cannot send negative amount",
		}

		c.IndentedJSON(http.StatusOK, res)
		return
	}

	res, err := r.todoService.Withdraw(postData.Amount)
	if err != nil {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(http.StatusOK, res)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.HttpResponse{
		Message: "success",
		Status:  http.StatusOK,
		Data:    res,
	})

	return
}
func (r *Router) deposit(c *gin.Context) {
	var postData model.Transaction

	if err := c.BindJSON(&postData); err != nil {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(http.StatusOK, res)
		return
	}

	if postData.Amount < 0 {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    "Cannot send negative amount",
		}

		c.IndentedJSON(http.StatusOK, res)
		return
	}

	res, err := r.todoService.Deposit(postData.Amount)

	if err != nil {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(http.StatusOK, res)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.HttpResponse{
		Message: "success",
		Status:  http.StatusOK,
		Data:    res,
	})

	return

}
func (r *Router) transfer(c *gin.Context) {
	var postData model.Transaction

	if err := c.BindJSON(&postData); err != nil {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	if postData.To == "" {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    "To field is required",
		}

		c.IndentedJSON(res.Status, res)
		return
	}

	res, err := r.todoService.Transfer(postData.To, postData.Amount)

	if err != nil {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusBadRequest,
			Data:    err.Error(),
		}

		c.IndentedJSON(http.StatusOK, res)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.HttpResponse{
		Message: "success",
		Status:  http.StatusOK,
		Data:    res,
	})

	return

}
