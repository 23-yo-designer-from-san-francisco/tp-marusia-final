package main

import (
	"encoding/json"
	"fmt"
	"guessTheSongMarusia/answer"
	"guessTheSongMarusia/game"
	"guessTheSongMarusia/models"
	"guessTheSongMarusia/router"
	"guessTheSongMarusia/utils"
	"math/rand"
	"os"
	"strings"
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
	musicR := musicRepo.NewMusicRepository(postgresDB)
	musicU := musicUsecase.NewMusicUsecase(musicR)
	musicD := musicDelivery.NewMusicDelivery(musicU)

	mywh := marusia.NewWebhook()
	mywh.EnableDebuging()

	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	b, err := os.ReadFile(`../../config/music.json`)
	if err != nil {
		log.Error(err)
	}
	jsonTracks := string(b)
	var allTracks models.TracksPerGenres
	if err := json.Unmarshal([]byte(jsonTracks), &allTracks); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	trackCount, err := musicU.GetTracksCount()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Warnf("Track count %d", trackCount)
	sessions := make(map[string]*models.Session)

	mywh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		userSession, ok := sessions[r.Session.SessionID]
		if !ok {
			userSession = models.NewSession()
			sessions[r.Session.SessionID] = userSession
			resp.Text, resp.TTS = answer.StartGamePhrase()
			return resp
		}

		switch r.Request.Type {
		case marusia.SimpleUtterance:
			// выход из игры в любом месте
			if r.Request.Command == marusia.OnInterrupt {
				resp.Text, resp.TTS = answer.GoodBye, answer.GoodBye
				resp.EndSession = true
				//TODO перенести сессии в базку или редиску
				delete(sessions, r.Session.SessionID)
				return resp
			}

			// TODO вместо strings.ContainsAny проверять наличие в токенах
			logrus.Warnf("Current mode: %d", userSession.GameStatus)
			switch userSession.GameStatus {
			case models.New:
				// логика после приветствия
				if utils.ContainsAny(r.Request.Command, answer.AgainE, answer.DontUnderstand, answer.Again) {
					// TODO хорошо бы состояние сделать "классом" со своей стандартной фразой, и запускать повторение прям из логики класса состояния (повторение будет перезапускать стандартную фразу состояния)
					resp.Text, resp.TTS = answer.StartGamePhrase()
				} else if strings.Contains(r.Request.Command, answer.Competition) {
					userSession.GameStatus = models.CompetitionIntro
					resp.Text, resp.TTS = answer.CompetitionRules, answer.CompetitionRules
				} else {
					userSession.GameStatus = models.ChoosingGenre
					resp.Text, resp.TTS = answer.ChooseGenre, answer.ChooseGenre
				}
			case models.ChoosingGenre, models.ListingGenres:
				// логика после предложения выбрать жанр
				if utils.ContainsAny(r.Request.Command, answer.AgainE, answer.DontUnderstand, answer.Again) {
					// попросили повторить
					switch userSession.GameStatus {
					case models.ChoosingGenre:
						resp.Text, resp.TTS = answer.ChooseGenre, answer.ChooseGenre
					case models.ListingGenres:
						resp.Text, resp.TTS = answer.AvailableGenres, answer.AvailableGenres
					}
				} else if utils.ContainsAny(r.Request.Command, answer.List, answer.LetsGo, answer.Available) {
					// попросили перечислить
					userSession.GameStatus = models.ListingGenres
					resp.Text, resp.TTS = answer.AvailableGenres, answer.AvailableGenres
				} else {
					resp = game.SelectGenre(userSession, r.Request.Command, resp, trackCount, musicU, sessions, allTracks, r.Session.SessionID, rng)
				}
			case models.Playing:
				// логика во время игры
				if utils.ContainsAny(r.Request.Command, answer.ChangeGenre, answer.ChangeGenre_, answer.AnotherGenre) {
					// попросили поменять жанр
					userSession.GameStatus = models.ChoosingGenre
					resp.Text, resp.TTS = answer.ChooseGenre, answer.ChooseGenre
				} else if userSession.MusicStarted {
					// после первого прослушивания
					if utils.ContainsAny(r.Request.Command, answer.Next, answer.GiveUp) {
						// игрок сдается
						userSession.MusicStarted = false
						resp.Text, resp.TTS = answer.LosePhrase(userSession)
					} else if utils.ContainsAll(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Title),
						strings.ToLower(userSession.CurrentTrack.Artist)) {
						// если сразу угадал исполнителя и название
						resp.Text, resp.TTS = answer.WinPhrase(userSession)
						userSession.MusicStarted = false
					} else if strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Title)) {
						// если угадал название
						userSession.TitleMatch = true
						if userSession.ArtistMatch {
							// если до этого угадал исполнителя
							resp.Text, resp.TTS = answer.WinPhrase(userSession)
							userSession.MusicStarted = false
						} else {
							resp = game.WrongAnswerPlay(userSession, resp)
						}
					} else if strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Artist)) {
						// если угадал исполнителя
						userSession.ArtistMatch = true
						if userSession.TitleMatch {
							// если до этого угадал название
							resp.Text, resp.TTS = answer.WinPhrase(userSession)
							userSession.MusicStarted = false
						} else {
							resp = game.WrongAnswerPlay(userSession, resp)
						}
					} else if userSession.NextLevelLoses {
						// если все попытки провалились
						userSession.MusicStarted = false
						resp.Text, resp.TTS = answer.LosePhrase(userSession)
					} else {
						resp = game.WrongAnswerPlay(userSession, resp)
					}
				} else {
					// перед первым или после последнего прослушивания
					resp = game.StartGame(userSession, resp, trackCount, musicU, rng)
				}
			case models.CompetitionIntro:
				if strings.Contains(r.Request.Command, answer.Competition) {
					userSession.GameStatus = models.CompetitionRules
					sessions[r.Session.SessionID] = userSession //TODO я забыл как делать нормально и хочу немного поспать
					resp.Text, resp.TTS = answer.CompetitionRules, answer.CompetitionRules
				}
			case models.CompetitionRules:
				if strings.Contains(r.Request.Command, answer.LetsGo) {
					logrus.Warn("COMPETITION MODE")
					userSession.GameStatus = models.Competition
				}
			case models.Competition:
				userSession.GameStatus = models.Playing
				resp.TTS = "скажите играем"
				resp.Text = resp.TTS
			default:
				resp.Text, resp.TTS = answer.IDontUnderstandYouPhrase()
			}
		}
		logrus.Warnf("New mode: %d", userSession.GameStatus)
		return resp
	})

	r := gin.Default()
	r.Any("/", gin.WrapF(mywh.HandleFunc))
	musicRouter := r.Group("/music")
	router.MusicEndpoints(musicRouter, musicD)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
