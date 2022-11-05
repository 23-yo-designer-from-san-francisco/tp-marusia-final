package main

import (
	"encoding/json"
	"fmt"
	"github.com/seehuhn/mt19937"
	"guessTheSongMarusia/answer"
	"guessTheSongMarusia/game"
	"guessTheSongMarusia/models"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/marusia"
)

// Навык "Угадай музло"
func main() {
	mywh := marusia.NewWebhook()
	mywh.EnableDebuging()

	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	b, err := os.ReadFile(`/Users/frbgd/sources/tp/tp-marusia-final/cmd/music.json`)
	if err != nil {
		fmt.Print(err)
	}
	jsonTracks := string(b)
	var allTracks models.TracksPerGenres
	if err := json.Unmarshal([]byte(jsonTracks), &allTracks); err != nil {
		log.Fatalln(err.Error())
	}
	var currentGameTracks []models.VKTrack

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
			fmt.Println("ok: ", ok, "music started: ", userSession.MusicStarted)
			fmt.Println("Command: ", r.Request.Command, "Track name: ", userSession.CurrentTrack.Title)

			// выход из игры в любом месте
			if r.Request.Command == marusia.OnInterrupt {
				resp.Text, resp.TTS = answer.GoodBye, answer.GoodBye
				resp.EndSession = true
				delete(sessions, r.Session.SessionID)
				return resp
			}

			// TODO вместо strings.Contains проверять наличие в токенах

			if userSession.GameStatus == models.New {
				// логика после приветствия
				if strings.Contains(r.Request.Command, answer.AgainE) || strings.Contains(r.Request.Command, answer.DontUnderstand) || strings.Contains(r.Request.Command, answer.Again) {
					// TODO хорошо бы состояние сделать "классом" со своей стандартной фразой, и запускать повторение прям из логики класса состояния (повторение будет перезапускать стандартную фразу состояния)
					resp.Text, resp.TTS = answer.StartGamePhrase()
				} else {
					userSession.GameStatus = models.ChoosingGenre
					resp.Text, resp.TTS = answer.ChooseGenre, answer.ChooseGenre
				}
			} else if userSession.GameStatus == models.ChoosingGenre || userSession.GameStatus == models.ListingGenres {
				// логика после предложения выбрать жанр|
				if strings.Contains(r.Request.Command, answer.AgainE) || strings.Contains(r.Request.Command, answer.DontUnderstand) || strings.Contains(r.Request.Command, answer.Again) {
					// попросили повторить
					if userSession.GameStatus == models.ChoosingGenre {
						resp.Text, resp.TTS = answer.ChooseGenre, answer.ChooseGenre
					} else if userSession.GameStatus == models.ListingGenres {
						resp.Text, resp.TTS = answer.AvailableGenres, answer.AvailableGenres
					}
				} else if strings.Contains(r.Request.Command, answer.List) || strings.Contains(r.Request.Command, answer.LetsGo) || strings.Contains(r.Request.Command, answer.Available) {
					// попросили перечислить
					userSession.GameStatus = models.ListingGenres
					resp.Text, resp.TTS = answer.AvailableGenres, answer.AvailableGenres
				} else if strings.Contains(r.Request.Command, strings.ToLower(answer.NotRock)) {
					// не рок
					// TODO вставить фразу о запуске не рока
					currentGameTracks = allTracks.NotRock
					sessions[r.Session.SessionID] = userSession
					userSession.GenreTrackCounter = 0
					userSession.CurrentGenre = answer.NotRock
					resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				} else if strings.Contains(r.Request.Command, strings.ToLower(answer.Rock)) {
					// рок
					// TODO вставить фразу о запуске рока
					currentGameTracks = allTracks.Rock
					sessions[r.Session.SessionID] = userSession
					userSession.GenreTrackCounter = 0
					userSession.CurrentGenre = answer.Rock
					resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				} else if strings.Contains(r.Request.Command, strings.ToLower(answer.Any)) {
					// любой
					// TODO вставить фразу о запуске любого
					currentGameTracks = append(allTracks.NotRock, allTracks.Rock...)
					sessions[r.Session.SessionID] = userSession
					userSession.GenreTrackCounter = 0
					userSession.CurrentGenre = answer.Any
					resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				} else {
					// непонел
					// TODO здесь надо находить жанр, похожий на названный
					resp.Text, resp.TTS = answer.IDontUnderstandYouPhrase()
				}
			} else if userSession.GameStatus == models.Playing {
				// логика во время игры
				if strings.Contains(r.Request.Command, answer.ChangeGenre) || strings.Contains(r.Request.Command, answer.ChangeGenre_) || strings.Contains(r.Request.Command, answer.AnotherGenre) {
					// попросили поменять жанр
					userSession.GameStatus = models.ChoosingGenre
					resp.Text, resp.TTS = answer.ChooseGenre, answer.ChooseGenre
				} else if userSession.MusicStarted {
					// после первого прослушивания
					if strings.Contains(r.Request.Command, answer.Next) || strings.Contains(r.Request.Command, answer.GiveUp) {
						// игрок сдается
						userSession.MusicStarted = false
						resp.Text, resp.TTS = answer.LosePhrase(userSession)
					} else if strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Title)) && strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Artist)) {
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
					resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				}
			} else {
				resp.Text, resp.TTS = answer.IDontUnderstandYouPhrase()
			}
		}
		return resp
	})

	http.HandleFunc("/", mywh.HandleFunc)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
