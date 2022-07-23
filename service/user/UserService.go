package service

import (
	"github.com/Tushar456/go-gin/entity"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Save(ctx *gin.Context, user *entity.User) error
	Get(ctx *gin.Context, user *string) (*entity.User, error)
	GetAll(ctx *gin.Context) ([]*entity.User, error)
	Update(ctx *gin.Context, user *entity.User) error
	Delete(ctx *gin.Context, username *string) error
}
