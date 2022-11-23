package router

import (
	musicD "guessTheSongMarusia/microservice/music/delivery"
	playlistDelivery "guessTheSongMarusia/microservice/playlist/delivery"

	"github.com/gin-gonic/gin"
)

func MusicEndpoints(r *gin.RouterGroup, mD *musicD.MusicDelivery) {
	r.POST("", mD.CreateAllMusic)
	r.GET("", mD.GetTracks)
}

func PlaylistEndpoints(r *gin.RouterGroup, pD *playlistDelivery.PlaylistDelivery) {
	r.POST("/create", pD.CreatePlaylist)
	r.GET("/generate-title", pD.GenerateTitle)
}
