package main

import (
	"fmt"
	"guessTheSongMarusia/game"
	"guessTheSongMarusia/middleware"
	"guessTheSongMarusia/router"
	"guessTheSongMarusia/utils"
	"math/rand"
	"os"
	"time"

	log "guessTheSongMarusia/pkg/logger"

	"github.com/seehuhn/mt19937"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/SevereCloud/vksdk/v2/marusia"
	"github.com/gin-gonic/gin"

	musicDelivery "guessTheSongMarusia/microservice/music/delivery"
	musicRepo "guessTheSongMarusia/microservice/music/repository"
	musicUsecase "guessTheSongMarusia/microservice/music/usecase"

	sessionRepo "guessTheSongMarusia/microservice/user/repository"
	sessionUsecase "guessTheSongMarusia/microservice/user/usecase"

	playlistDelivery "guessTheSongMarusia/microservice/playlist/delivery"
	playlistRepo "guessTheSongMarusia/microservice/playlist/repository"
	playlistUsecase "guessTheSongMarusia/microservice/playlist/usecase"
)

const logMessage = "server:"

// Навык "Угадай музло"
func main() {
	message := logMessage + "Main:"
	log.Init(logrus.DebugLevel)
	log.Info(fmt.Sprintf(message+"started, log level = %s", logrus.DebugLevel))

	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(".env file isn't found")
		os.Exit(1)
	}

	postgresDB, err := utils.InitPostgres()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	redisDB, err := utils.InitRedisDB()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	adjectives := utils.ReadCsvFile(viper.GetString("ADJECTIVES_FILEPATH"))
	nouns := utils.ReadCsvFile(viper.GetString("NOUNS_FILEPATH"))

	middleware := middleware.NewMiddleware()

	musicR := musicRepo.NewMusicRepository(postgresDB)
	musicU := musicUsecase.NewMusicUsecase(musicR)
	musicD := musicDelivery.NewMusicDelivery(musicU)

	playlistR := playlistRepo.NewPlaylistRepository(redisDB)
	playlistU := playlistUsecase.NewPlaylistUsecase(adjectives, nouns, playlistR, musicR)
	playlistD := playlistDelivery.NewPlaylistDelivery(playlistU)

	sessionR := sessionRepo.NewSessionRepository(redisDB)
	sessionU := sessionUsecase.NewUserUsecase(sessionR)

	mywh := marusia.NewWebhook()
	mywh.EnableDebuging()

	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	mywh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		log.Debug("Got request command:", r.Request.Command)
		return game.MainHandler(r, sessionU, musicU, playlistU, rng, adjectives, nouns)
	})

	r := gin.Default()
	r.POST("/", gin.WrapF(mywh.HandleFunc))
	r.OPTIONS("/", gin.WrapF(mywh.HandleFunc))

	musicRouter := r.Group("/music")
	musicRouter.Use(middleware.CORSMiddleware())

	playlistRouter := r.Group("/playlists")
	playlistRouter.Use(middleware.CORSMiddleware())

	router.MusicEndpoints(musicRouter, musicD)
	router.PlaylistEndpoints(playlistRouter, playlistD)

	err = r.Run(fmt.Sprintf(":%d", viper.GetInt("SERVER_PORT")))
	if err != nil {
		return
	}
}
