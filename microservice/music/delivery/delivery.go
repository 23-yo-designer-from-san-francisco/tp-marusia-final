package delivery

import (
	//"encoding/json"
	"guessTheSongMarusia/microservice/music"
	"guessTheSongMarusia/models"
	"net/http"

	log "guessTheSongMarusia/pkg/logger"

	"github.com/gin-gonic/gin"
)

const logMessage = "microservice:music:delivery:"

type MusicDelivery struct {
	musicUsecase music.Usecase
}

func NewMusicDelivery(musicU music.Usecase) *MusicDelivery {
	return &MusicDelivery{
		musicUsecase: musicU,
	}
}

func (mD *MusicDelivery) CreateAllMusic(c *gin.Context) {
	message := logMessage + "CreateAllMusic:"
	log.Debug(message + "started")

	vkTracks := []models.VKTrack{}
	err := c.ShouldBindJSON(&vkTracks)
	if err != nil {
		log.Error(message+"err = ", err)
		c.JSON(http.StatusOK, err.Error())
		return
	}
	err = mD.musicUsecase.CreateAllMusic(vkTracks)
	if err != nil {
		log.Error(message+"err = ", err)
		c.JSON(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func (mD *MusicDelivery) GetTracks(c *gin.Context) {
	message := logMessage + "GetTracks:"
	log.Debug(message + "started")

	var resultSongs []models.VKTrack
	var err error

	artist := c.Query("artist")
	if artist == "" {
		resultSongs, err = mD.musicUsecase.GetAllMusic()
	} else {
		resultSongs, _, err = mD.musicUsecase.GetSongsByArtist(artist)
	}

	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	var miniAppTracks []models.MiniAppTrack
	for _, song := range resultSongs {
		miniAppTracks = append(miniAppTracks, models.MiniAppTrack{
			ID: song.ID,
			Title: song.Title,
			Artist: song.Artist,
		})
	}
	c.JSON(http.StatusOK, miniAppTracks)
}
