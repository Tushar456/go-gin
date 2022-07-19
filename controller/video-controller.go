package controller

import (
	"net/http"

	"github.com/Tushar456/go-gin/entity"
	"github.com/Tushar456/go-gin/service"
	"github.com/gin-gonic/gin"
)

type Videocontroller struct {
	service service.VideoService
}

func (controller *Videocontroller) GetALLVideos(ctx *gin.Context) {
	videos := controller.service.FindAll()
	ctx.JSON(http.StatusOK, videos)

}

func (controller *Videocontroller) CreateVideo(ctx *gin.Context) {

	var video entity.Video

	if err := ctx.BindJSON(&video); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	controller.service.Save(video)
	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})

}

func (vc *Videocontroller) RegisterVideoRoutes(rg *gin.RouterGroup) {
	videoroute := rg.Group("/video")
	videoroute.POST("", vc.CreateVideo)
	videoroute.GET("", vc.GetALLVideos)

}

func New(service service.VideoService) Videocontroller {
	return Videocontroller{
		service: service,
	}
}
