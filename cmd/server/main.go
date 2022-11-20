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
	nouns := utils.ReadCsvFile("../../config/output.txt")

	musicR := musicRepo.NewMusicRepository(postgresDB)
	musicU := musicUsecase.NewMusicUsecase(musicR)
	musicD := musicDelivery.NewMusicDelivery(musicU)

	sessionR := sessionRepo.NewSessionRepository(redisDB)
	sessionU := sessionUsecase.NewUserUsecase(sessionR)

	mywh := marusia.NewWebhook()
	mywh.EnableDebuging()

	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	mywh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		return game.MainHandler(r, sessionU, musicU, rng, adjectives, nouns)
	})

	r := gin.New()
	r.Use(gin.Recovery())
	r.Any("/", gin.WrapF(mywh.HandleFunc))
	musicRouter := r.Group("/music")
	router.MusicEndpoints(musicRouter, musicD)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
