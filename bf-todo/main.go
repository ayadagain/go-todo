package main

import (
	"assm/bf-todo/ctx"
	"assm/bf-todo/http/route"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")
	serviceContext := ctx.NewDefaultServiceContext()
	router := route.NewRouter(serviceContext)
	ginEngine := gin.Default(router.Install)
	err := ginEngine.Run()
	if err != nil {
		panic(err)
	}
}
