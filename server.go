package main

import (
	"github.com/Tushar456/go-gin/controller"
	ucontoller "github.com/Tushar456/go-gin/controller/user"
	"github.com/Tushar456/go-gin/middleware"
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
	jwtToken, err := token.NewJWTToken(tokenSymmetricKey)
	if err != nil {
		return
	}
	router := gin.Default()
	baseroute := router.Group("api/v1/user")
	baseroute.POST("/create", userController.CreateUser)
	baseroute.POST("/login", userController.LoginUser)
	userroute := baseroute.Group("").Use(middleware.AuthMiddleware(jwtToken))

	userroute.GET("/:username", userController.GetUser)
	userroute.GET("", userController.GetAll)
	userroute.DELETE("/:username", userController.DeleteUser)
	userroute.PUT("", userController.UpdateUser)

	router.Run(":8080")

}

func pingHandler(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"message": "pong is working",
	})

}
