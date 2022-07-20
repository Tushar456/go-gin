package controller

import (
	"errors"
	"github.com/Tushar456/go-gin/entity"
	"github.com/Tushar456/go-gin/helper"
	"github.com/Tushar456/go-gin/middleware"
	service "github.com/Tushar456/go-gin/service/user"
	"github.com/Tushar456/go-gin/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

var validateError validator.ValidationErrors
var jwtToken token.Token

type UserController struct {
	UserService service.UserService
}

func New(userservice service.UserService) UserController {
	return UserController{
		UserService: userservice,
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user entity.User
	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {

		if errors.As(err, &validateError) {
			helper.CustomerValidateErrorMessage(ctx, validateError)
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.UserService.Save(&user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	var username string = ctx.Param("username")
	user, err := uc.UserService.Get(&username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAll(ctx *gin.Context) {
	users, err := uc.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user entity.User
	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {

		if errors.As(err, &validateError) {
			helper.CustomerValidateErrorMessage(ctx, validateError)
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := uc.UserService.Update(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	var username string = ctx.Param("username")
	err := uc.UserService.Delete(&username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

type loginUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := uc.UserService.Get(&req.UserName)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if user.Password != req.Password {

		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("password is incorrect")))
		return

	}
	accessToken, accessPayload, err := jwtToken.CreateToken(user.UserName, time.Duration(1)*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup, token token.Token) {
	jwtToken = token
	userroute := rg.Group("/user")
	userroute.POST("/create", uc.CreateUser)
	userroute.POST("/login", uc.LoginUser)

	userroute.GET("/:username", uc.GetUser).Use(middleware.AuthMiddleware(jwtToken))
	userroute.GET("", uc.GetAll).Use(middleware.AuthMiddleware(jwtToken))
	userroute.DELETE("/:username", uc.DeleteUser).Use(middleware.AuthMiddleware(jwtToken))
	userroute.PUT("", uc.UpdateUser).Use(middleware.AuthMiddleware(jwtToken))
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
