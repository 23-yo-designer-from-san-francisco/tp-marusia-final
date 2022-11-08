package router

import (
	musicD "guessTheSongMarusia/microservice/music/delivery"

	"github.com/gin-gonic/gin"
)

func MusicEndpoints(r *gin.RouterGroup, mD *musicD.MusicDelivery) {
	r.POST("", mD.CreateAllMusic)
}