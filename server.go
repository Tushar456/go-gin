package main

import (
	"github.com/Tushar456/go-gin/controller"
	ucontoller "github.com/Tushar456/go-gin/controller/user"
	"github.com/Tushar456/go-gin/service"
	uservice "github.com/Tushar456/go-gin/service/user"
	"github.com/Tushar456/go-gin/token"

	"github.com/gin-gonic/gin"
)

var (
	videoService    service.VideoService       = service.New()
	videoController controller.Videocontroller = controller.New(videoService)
	userService     uservice.UserService       = uservice.NewUserService()
	userController  ucontoller.UserController  = ucontoller.New(userService)
)

const (
	tokenSymmetricKey = "12345678901234567890123456789012"
)

func main() {
	token, err := token.NewJWTToken(tokenSymmetricKey)
	if err != nil {
		return
	}
	server := gin.Default()
	baspath := server.Group("/api/v1")
	baspath.GET("/ping", pingHandler)

	videoController.RegisterVideoRoutes(baspath)
	userController.RegisterUserRoutes(baspath, token)
	server.Run(":8080")

}

func pingHandler(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"message": "pong is working",
	})

}
