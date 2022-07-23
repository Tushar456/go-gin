package controller

import (
	"errors"
	"github.com/Tushar456/go-gin/entity"
	"github.com/Tushar456/go-gin/helper"
	service "github.com/Tushar456/go-gin/service/user"
	"github.com/Tushar456/go-gin/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

const (
	authorizationPayloadKey = "payload"
	tokenSymmetricKey       = "12345678901234567890123456789012"
)

var validateError validator.ValidationErrors

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

	if err := uc.UserService.Save(ctx, &user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	var username string = ctx.Param("username")
	authPayload := ctx.MustGet(authorizationPayloadKey)
	payload := authPayload.(*token.Payload)
	if payload.Username != username && payload.Username != "admin" {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("user doesnot have write to access")))
		return
	}
	user, err := uc.UserService.Get(ctx, &username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAll(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey)
	payload := authPayload.(*token.Payload)
	if payload.Username != "admin" {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("user doesnot have write to access")))
		return
	}

	users, err := uc.UserService.GetAll(ctx)
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
	authPayload := ctx.MustGet(authorizationPayloadKey)
	payload := authPayload.(*token.Payload)
	if payload.Username != user.UserName && payload.Username != "admin" {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("user doesnot have write to access")))
		return
	}

	err := uc.UserService.Update(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	var username string = ctx.Param("username")

	authPayload := ctx.MustGet(authorizationPayloadKey)
	payload := authPayload.(*token.Payload)
	if payload.Username != username && payload.Username != "admin" {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("user doesnot have write to access")))
		return
	}
	err := uc.UserService.Delete(ctx, &username)
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
	jwtToken, err := token.NewJWTToken(tokenSymmetricKey)
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := uc.UserService.Get(ctx, &req.UserName)
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

//func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup, token token.Token) {
//	jwtToken = token
//	userroute := rg.Group("/user")
//	userroute.GET("/:username", uc.GetUser).Use(middleware.AuthMiddleware(jwtToken))
//	userroute.GET("", uc.GetAll).Use(middleware.AuthMiddleware(jwtToken))
//	userroute.DELETE("/:username", uc.DeleteUser).Use(middleware.AuthMiddleware(jwtToken))
//	userroute.PUT("", uc.UpdateUser).Use(middleware.AuthMiddleware(jwtToken))
//	userroute.POST("/create", uc.CreateUser)
//	userroute.POST("/login", uc.LoginUser)
//}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
