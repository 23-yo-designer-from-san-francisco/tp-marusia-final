package delivery

import (
	"guessTheSongMarusia/microservice/playlist"
	"guessTheSongMarusia/models"
	"net/http"

	"github.com/gin-gonic/gin"
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
	var playlist models.MiniAppPlaylist
	err := c.ShouldBindJSON(&playlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = pD.playlistUsecase.CreatePlaylist(playlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func (pD *PlaylistDelivery) GenerateTitle(c *gin.Context) {
	titleKey, title, err := pD.playlistUsecase.GenerateTitle()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.TitleInfo{
		TitleKey: titleKey,
		Title: title,
	})
}