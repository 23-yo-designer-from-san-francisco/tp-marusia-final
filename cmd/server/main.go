package main

import (
	"fmt"
	"guessTheSongMarusia/game"
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

	viper.AddConfigPath("../../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Config isn't found 1")
		os.Exit(1)
	}
	viper.SetConfigFile("../../.env")
	err = viper.MergeInConfig()
	if err != nil {
		log.Error("Config isn't found 2")
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

	adjectives := utils.ReadCsvFile("../../config/adjectives.txt")
	nouns := utils.ReadCsvFile("../../config/nouns.txt")

	musicR := musicRepo.NewMusicRepository(postgresDB)
	musicU := musicUsecase.NewMusicUsecase(musicR)
	musicD := musicDelivery.NewMusicDelivery(musicU)

	playlistR := playlistRepo.NewPlaylistRepository(redisDB)
	playlistU := playlistUsecase.NewPlaylistUsecase(playlistR)
	playlistD := playlistDelivery.NewPlaylistDelivery(playlistU)

	sessionR := sessionRepo.NewSessionRepository(redisDB)
	sessionU := sessionUsecase.NewUserUsecase(sessionR)

	mywh := marusia.NewWebhook()
	mywh.EnableDebuging()

	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	mywh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		return game.MainHandler(r, sessionU, musicU, playlistU, rng, adjectives, nouns)
	})

	r := gin.New()
	r.Use(gin.Recovery())
	r.Any("/", gin.WrapF(mywh.HandleFunc))
	musicRouter := r.Group("/music")
	playlistRouter := r.Group("/playlists")
	router.MusicEndpoints(musicRouter, musicD)
	router.PlaylistEndpoints(playlistRouter, playlistD)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
