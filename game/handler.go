package game

import (
	"fmt"
	musicUsecase "guessTheSongMarusia/microservice/music/usecase"
	playlistUsecase "guessTheSongMarusia/microservice/playlist/usecase"
	"guessTheSongMarusia/microservice/user"
	"guessTheSongMarusia/models"
	log "guessTheSongMarusia/pkg/logger"
	"guessTheSongMarusia/utils"
	"math/rand"
	"strings"

	"github.com/SevereCloud/vksdk/v2/marusia"
	"github.com/sirupsen/logrus"
)

func MainHandler(r marusia.Request,
	sessionU user.SessionUsecase, musicU *musicUsecase.MusicUsecase, playlistU *playlistUsecase.PlaylistUsecase,
	rng *rand.Rand, adjectives []string, nouns []string) (resp marusia.Response) {
	log.Debug("Got command:", r.Request.Command)
	userSession, err := sessionU.GetSession(r.Session.SessionID)
	if err != nil {
		log.Error(err.Error())
		userSession = models.NewSession()
		err := sessionU.SaveSession(r.Session.SessionID, userSession)
		if err != nil {
			log.Error(err.Error())
			resp.Text, resp.TTS = userSession.GameState.SayErrorPhrase()
			return resp
		}
		resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
		return resp
	}

	switch r.Request.Type {
	case marusia.SimpleUtterance:
		// выход из игры в любом месте
		if r.Request.Command == marusia.OnInterrupt {
			logrus.Debug("OnInterrupt", r, userSession)
			resp.Text, resp.TTS = models.GoodBye, models.GoodBye
			resp.EndSession = true
			err := sessionU.DeleteSession(r.Session.SessionID)
			if err != nil {
				log.Error(err.Error())
				resp.Text, resp.TTS = userSession.GameState.SayErrorPhrase()
				return resp
			}
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
		} else if utils.ContainsAny(
			r.Request.Command, models.ChooseArtist, models.ChangeArtist, models.ChangeArtist_,
			models.AnotherArtist, models.ChooseArtist__, models.ChangeArtist__, models.ChangeArtist___,
			models.AnotherArtist__,
		) {
			// попросили поменять артиста
			userSession.GameState = models.ChooseArtistState
			resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
		} else if userSession.CompetitionMode && utils.ContainsAny(r.Request.Command, models.CompetitionRule) {
			resp.Text, resp.TTS = models.CompetitionRulesPhrase()
		} else {
			switch userSession.GameState.GameStatus {
			case models.StatusNewGame:
				// логика после приветствия
				logrus.Debug("NewGameRequest", r, userSession)
				if strings.Contains(r.Request.Command, models.Competition) {
					userSession.GameState = models.CompetitonState
				} else if strings.Contains(r.Request.Command, models.LetsPlay) {
					userSession.GameState = models.ChooseGenreState
				} else if strings.Contains(r.Request.Command, models.KeyPhrase) {
					userSession.GameState = models.KeyPhrasePlaylistState
				}
				resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
				logrus.Debug("NewGameResponse", r, userSession)
			case models.StatusChoosingGenre, models.StatusListingGenres:
				// логика после предложения выбрать жанр
				logrus.Debug("GenresRequest", r, userSession)
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
				} else if utils.ContainsAny(r.Request.Command, models.Artist) {
					//Переходим на артистов
					userSession.GameState = models.ChooseArtistState
					resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
				} else {
					// ищем названный жанр и начинаем игру
					resp = SelectGenre(userSession, r.Request.Command, nouns, adjectives, resp, musicU, playlistU, rng)
				}
				logrus.Debug("GenresResponse", r, userSession)

			case models.StatusChooseArtist:
				logrus.Debug("ArtistRequest", r, userSession)
				if utils.ContainsAny(r.Request.Command, models.AgainE, models.DontUnderstand, models.Again) {
					// попросили повторить
					resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					logrus.Debug("ArtistRequest", r, userSession)
					return
				}
				if utils.ContainsAny(r.Request.Command, models.Genre) {
					// выход обратно к жанрам
					userSession.GameState = models.ChooseGenreState
					resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
				}
				// ищем названного исполнителя и начинаем игру
				resp = SelectArtist(userSession, r.Request.Command, nouns, adjectives, resp, musicU, playlistU, rng)
				logrus.Debug("ArtistRequest", r, userSession)

			case models.StatusPlaying:
				// логика во время игры
				logrus.Debug("PlayingRequest", r, userSession)
				if !userSession.MusicStarted {
					// перед первым или после последнего прослушивания (или если игра не стартанула из выбора жанра/артиста)
					resp = StartGame(userSession, resp)
				} else {
					// после первого прослушивания
					if utils.ContainsAny(r.Request.Command, models.Next, models.GiveUp) {
						logrus.Info("Gave up")
						// игрок сдается
						userSession.MusicStarted = false
						resp.Text, resp.TTS = models.LosePhrase(userSession)
					} else {
						matchTitle, matchArtists := userSession.CurrentTrack.CheckUserAnswer(r.Request.Command, userSession)
						log.Debug(matchTitle, matchArtists)
						if userSession.TitleMatch && userSession.ArtistMatch {
							logrus.Info("Guessed both")
							userSession.MusicStarted = false
							resp.Text, resp.TTS = models.WinPhrase(userSession)
						} else if matchTitle {
							logrus.Info("Guessed title")
							resp = CloseAnswerPlay(userSession, resp)
						} else if matchArtists {
							logrus.Info("Guessed artist")
							resp = CloseAnswerPlay(userSession, resp)
						} else if userSession.NextLevelLoses {
							// если все попытки провалились
							userSession.MusicStarted = false
							resp.Text, resp.TTS = models.LosePhrase(userSession)
						} else {
							resp = WrongAnswerPlay(userSession, resp)
						}
					}
				}
				logrus.Debug("PlayingResponse", r, userSession)
			case models.StatusCompetitionRules:
				if utils.ContainsAny(r.Request.Command, models.Yes, models.LetsPlay) {
					userSession.GameState = models.ChooseGenreState
					userSession.CompetitionMode = true
					userSession.CurrentPoints = 0
					resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
				} else if strings.Contains(r.Request.Command, models.No) {
					userSession.GameState = models.NewGameState
					resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
				}
				resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
			case models.StatusGeneratedPlaylist:
				userSession.GameState = models.NewGameState
				resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
			case models.StatusPlaylistFinished:
				resp.Text, resp.TTS = "", ""
				if userSession.CompetitionMode {
					resp.Text = models.GetScoreText(userSession)
					resp.TTS = models.GetScoreText(userSession)
				}
				standartText, standartTTS := userSession.GameState.SayStandartPhrase()
				resp.Text += standartText
				resp.TTS += standartTTS
			case models.StatusKeyPhrasePlaylist:
				playlist, err := playlistU.GetPlaylist(r.Request.Command)
				if err != nil {
					resp.Text, resp.TTS = userSession.GameState.SayErrorPhrase()
					return
				}
				length := len(playlist)
				if length == 0 {
					resp.Text, resp.TTS = userSession.GameState.SayStandartPhrase()
					return
				}
				userSession.CurrentGenre = r.Request.Command
				userSession.KeyPhrase = r.Request.Command
				userSession.CurrentPlaylist = playlist
				userSession.CompetitionMode = true
				str := fmt.Sprintf("%s%d ", "Я нашла этот плейлист. Количество треков: ", len(playlist))
				resp = StartGame(userSession, resp)
				resp.Text, resp.TTS = str+resp.Text, str+resp.TTS
			default:
				resp.Text, resp.TTS = models.IDontUnderstandYouPhrase()
			}
		}
	}
	logrus.Warnf("New mode: %d", userSession.GameState.GameStatus)
	err = sessionU.SaveSession(r.Session.SessionID, userSession)
	if err != nil {
		log.Error(err.Error())
		resp.Text, resp.TTS = userSession.GameState.SayErrorPhrase()
	}
	return resp
}
