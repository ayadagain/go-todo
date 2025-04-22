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

func newBadRequestResponse(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, &model.HttpResponse{
		Status:  http.StatusBadRequest,
		Message: message,
		Data:    nil,
	})
}

func newSuccessResponse(c *gin.Context, message string, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, &model.HttpResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func (r *Router) withdraw(c *gin.Context) {
	var postData model.Transaction

	if err := c.BindJSON(&postData); err != nil {
		newBadRequestResponse(c, err.Error())
		return
	}

	if postData.Amount < 0 {
		newBadRequestResponse(c, "Cannot send negative amount")
		return
	}

	res, err := r.todoService.Withdraw(postData.Amount)
	if err != nil {
		newBadRequestResponse(c, err.Error())
		return
	}

	newSuccessResponse(c, "success", res)
	return
}
func (r *Router) deposit(c *gin.Context) {
	var postData model.Transaction

	if err := c.BindJSON(&postData); err != nil {
		newBadRequestResponse(c, err.Error())
		return
	}

	if postData.Amount < 0 {
		newBadRequestResponse(c, "Cannot send negative amount")
		return
	}

	res, err := r.todoService.Deposit(postData.Amount)

	if err != nil {
		newBadRequestResponse(c, err.Error())
		return
	}

	newSuccessResponse(c, "success", res)
	return
}
func (r *Router) transfer(c *gin.Context) {
	var postData model.Transaction

	if err := c.BindJSON(&postData); err != nil {
		newBadRequestResponse(c, err.Error())
		return
	}

	if postData.To == "" {
		newBadRequestResponse(c, "To field is required")
		return
	}

	res, err := r.todoService.Transfer(postData.To, postData.Amount)

	if err != nil {
		newBadRequestResponse(c, err.Error())
		return
	}

	newSuccessResponse(c, "success", res)
	return
}
