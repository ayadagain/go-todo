package route

import (
	"assm/bf-todo/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *Router) ServerKeyCheck(c *gin.Context) {
	serverKey := c.Request.Header.Get("ServerKey")
	envServKey := r.serviceContext.Conf().GetSecretKey()

	if serverKey != envServKey {
		res := &model.HttpResponse{
			Message: "fail",
			Status:  http.StatusUnauthorized,
			Data:    nil,
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	c.Next()
}
