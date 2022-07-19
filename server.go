package main

import (
	"github.com/Tushar456/go-gin/controller"
	ucontoller "github.com/Tushar456/go-gin/controller/user"
	"github.com/Tushar456/go-gin/service"
	uservice "github.com/Tushar456/go-gin/service/user"

	"github.com/gin-gonic/gin"
)

var (
	videoService    service.VideoService       = service.New()
	videoController controller.Videocontroller = controller.New(videoService)
	userService     uservice.UserService       = uservice.NewUserService()
	userController  ucontoller.UserController  = ucontoller.New(userService)
)

func main() {

	server := gin.Default()
	baspath := server.Group("/api/v1")
	baspath.GET("/ping", pingHandler)

	videoController.RegisterVideoRoutes(baspath)
	userController.RegisterUserRoutes(baspath)
	server.Run(":8080")

}

func pingHandler(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"message": "pong is working",
	})

}
