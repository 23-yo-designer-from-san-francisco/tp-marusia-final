package main

import (
	"fmt"
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

func printLog(blockName string, r marusia.Request, userSession *models.Session) {
	logMessage := "Command: " + r.Request.Command + " "
	logMessage += "In " + blockName + " Block "
	logMessage += "SessionInfo: " + fmt.Sprint(*userSession) + " "
	logMessage += "GameStateInfo: " + fmt.Sprint(*userSession.GameState)
	//TODO логи выключил тут log.Debug(logMessage)
}

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

	sessions := make(map[string]*models.Session)

	mywh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		userSession, ok := sessions[r.Session.SessionID]
		if !ok {
			userSession = models.NewSession()
			sessions[r.Session.SessionID] = userSession
			resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
			return resp
		}

		switch r.Request.Type {
		case marusia.SimpleUtterance:
			// выход из игры в любом месте
			if r.Request.Command == marusia.OnInterrupt {
				printLog("OnInterrupt", r, userSession)
				resp.Text, resp.TTS = models.GoodBye, models.GoodBye
				resp.EndSession = true
				//TODO перенести сессии в базку или редиску
				delete(sessions, r.Session.SessionID)
				return resp
			}
			logrus.Warnf("Current mode: %d", userSession.GameState.GameStatus)
			if utils.ContainsAny(r.Request.Command, models.ChangeGame, models.ChangeGame_, models.AnotherGame) {
				// попросили поменять игру
				userSession.GameState = models.NewGameState
				resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
			} else if utils.ContainsAny(r.Request.Command, models.ChangeGenre, models.ChangeGenre_, models.AnotherGenre) {
				// попросили поменять жанр
				userSession.GameState = models.ChooseGenreState
				resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
			} else if utils.ContainsAny(r.Request.Command, models.ChooseArtist, models.ChangeArtist, models.ChangeArtist_, models.AnotherArtist, models.ChangeArtist__, models.ChangeArtist___, models.AnotherArtist__) {
				// попросили поменять артиста
				userSession.GameState = models.ChooseArtistState
				resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
			} else {
				switch userSession.GameState.GameStatus {
				case models.StatusNewGame:
					// логика после приветствия
					printLog("NewGameRequest", r, userSession)
					if strings.Contains(r.Request.Command, models.Competition) {
						userSession.GameState = models.CompetitonRulesState
						resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					} else if strings.Contains(r.Request.Command, models.LetsPlay) {
						userSession.GameState = models.ChooseGenreState
						resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					}
					printLog("NewGameResponse", r, userSession)
				case models.StatusChoosingGenre, models.StatusListingGenres:
					// логика после предложения выбрать жанр
					printLog("GenresRequest", r, userSession)
					if utils.ContainsAny(r.Request.Command, models.AgainE, models.DontUnderstand, models.Again) {
						// попросили повторить
						resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					} else if utils.ContainsAny(r.Request.Command, models.List, models.LetsGo, models.Available) {
						// попросили перечислить
						userSession.GameState = models.ListingGenreState
						genres, err := musicU.GetGenres()
						if err != nil {
							resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
						}
						str := "Предлагаю следующие жанры:\n"
						for _, genre := range genres {
							str += genre + "\n"
						}
						resp.Text, resp.TTS = str, str
					} else if utils.ContainsAny(r.Request.Command, models.Artists) {
						//Переходим на артистов
						userSession.GameState = models.ChooseArtistState
						resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					} else {
						// ищем названный жанр и начинаем игру
						resp = game.SelectGenre(userSession, r.Request.Command, resp, musicU, sessions, r.Session.SessionID, rng)
					}
					printLog("GenresResponse", r, userSession)

				case models.StatusChooseArtist:
					printLog("ArtistRequest", r, userSession)
					if utils.ContainsAny(r.Request.Command, models.AgainE, models.DontUnderstand, models.Again) {
						// попросили повторить
						resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
						printLog("ArtistRequest", r, userSession)
						return
					}
					if utils.ContainsAny(r.Request.Command, models.Genre) {
						// выход обратно к жанрам
						userSession.GameState = models.ChooseGenreState
						resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					}
					// ищем названного исполнителя и начинаем игру
					resp = game.SelectArtist(userSession, r.Request.Command, resp, musicU, sessions, r.Session.SessionID, rng)
					printLog("ArtistRequest", r, userSession)

				case models.StatusPlaying:
					// логика во время игры
					printLog("PlayingRequest", r, userSession)
					if !userSession.MusicStarted {
						// перед первым или после последнего прослушивания (или если игра не стартанула из выбора жанра/артиста)
						resp = game.StartGame(userSession, resp, musicU, rng)
					} else {
						// после первого прослушивания
						if utils.ContainsAny(r.Request.Command, models.Next, models.GiveUp) {
							logrus.Info("Gave up")
							// игрок сдается
							userSession.MusicStarted = false
							resp.Text, resp.TTS = models.LosePhrase(userSession)
						} else {
							matchTitle, matchArtists := userSession.CurrentTrack.CheckUserAnswer(r.Request.Command)
							if !userSession.TitleMatch && !userSession.ArtistMatch && matchTitle && matchArtists {
								logrus.Info("Guessed both")
								resp.Text, resp.TTS = models.WinPhrase(userSession)
								switch userSession.CurrentLevel {
								case models.Two:
									userSession.CurrentPoints += models.GuessedAttempt1
								case models.Five:
									userSession.CurrentPoints += models.GuessedAttempt2
								case models.Ten:
									userSession.CurrentPoints += models.GuessedAttempt3
								}
								userSession.MusicStarted = false
							} else if !userSession.TitleMatch && matchTitle {
								logrus.Info("Guessed title")
								userSession.TitleMatch = true
								switch userSession.CurrentLevel {
								case models.Two:
									userSession.CurrentPoints += models.GuessedAttempt1 / 2
								case models.Five:
									userSession.CurrentPoints += models.GuessedAttempt2 / 2
								case models.Ten:
									userSession.CurrentPoints += models.GuessedAttempt3 / 2
								}
								if userSession.ArtistMatch || userSession.GameMode == models.ArtistMode {
									// если до этого угадал исполнителя
									logrus.Info("Guessed artist before and title now")
									resp.Text, resp.TTS = models.WinPhrase(userSession)
									userSession.MusicStarted = false
								} else {
									resp = game.CloseAnswerPlay(userSession, resp)
								}
							} else if !userSession.ArtistMatch && matchArtists {
								logrus.Info("Guessed artist")
								userSession.ArtistMatch = true
								switch userSession.CurrentLevel {
								case models.Two:
									userSession.CurrentPoints += models.GuessedAttempt1 / 2
								case models.Five:
									userSession.CurrentPoints += models.GuessedAttempt2 / 2
								case models.Ten:
									userSession.CurrentPoints += models.GuessedAttempt3 / 2
								}
								if userSession.TitleMatch {
									// если до этого угадал название
									logrus.Warn("Guessed title before and artist now")
									resp.Text, resp.TTS = models.WinPhrase(userSession)
									userSession.MusicStarted = false
								} else {
									resp = game.CloseAnswerPlay(userSession, resp)
								}
							} else if userSession.NextLevelLoses {
								// если все попытки провалились
								userSession.MusicStarted = false
								resp.Text, resp.TTS = models.LosePhrase(userSession)
							} else {
								resp = game.WrongAnswerPlay(userSession, resp)
							}
						}
					}
					printLog("PlayingResponse", r, userSession)
				case models.StatusCompetitionRules:
					if strings.Contains(r.Request.Command, models.Yes) {
						userSession.GameState = models.ChooseGenreState
						userSession.CompetitionMode = true
						userSession.CurrentPoints = 0
						resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					} else if strings.Contains(r.Request.Command, models.No) {
						userSession.GameState = models.NewGameState
						resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					}
				default:
					resp.Text, resp.TTS = models.IDontUnderstandYouPhrase()
				}
			}
		}
		logrus.Warnf("New mode: %d", userSession.GameState.GameStatus)
		return resp
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
