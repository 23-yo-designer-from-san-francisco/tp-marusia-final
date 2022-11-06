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
	//err := json.NewDecoder(c.Request.Body).Decode(&vkTracks)
	err := c.ShouldBindJSON(&vkTracks)
	if err != nil {
		log.Error(message + "err = ", err)
		c.JSON(http.StatusOK, err.Error())
		return
	}
	err = mD.musicUsecase.CreateAllMusic(vkTracks)
	if err != nil {
		log.Error(message + "err = ", err)
		c.JSON(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, "OK")
}