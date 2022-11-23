package delivery

import (
	"github.com/gin-gonic/gin"
	"guessTheSongMarusia/microservice/playlist"
)

type PlaylistDelivery struct {
	playlistUsecase playlist.Usecase
}

func NewPlaylistDelivery(playlistU playlist.Usecase) *PlaylistDelivery {
	return &PlaylistDelivery{
		playlistUsecase: playlistU,
	}
}

func (pD *PlaylistDelivery) CreatePlaylist(c *gin.Context) {

}

func (pD *PlaylistDelivery) GenerateTitle(c *gin.Context) {

}
