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
					if userSession.GameStatus == models.ChoosingGenre {
						resp.Text, resp.TTS = answer.ChooseGenre, answer.ChooseGenre
					} else if userSession.GameStatus == models.ListingGenres {
						resp.Text, resp.TTS = answer.AvailableGenres, answer.AvailableGenres
					}
				} else if strings.Contains(r.Request.Command, answer.List) || strings.Contains(r.Request.Command, answer.Available) {
					userSession.GameStatus = models.ListingGenres
					resp.Text, resp.TTS = answer.AvailableGenres, answer.AvailableGenres
				} else if strings.Contains(r.Request.Command, answer.NotRock) {
					// TODO вставить фразу о запуске не рока
					currentGameTracks = allTracks.NotRock
					sessions[r.Session.SessionID] = userSession
					resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				} else if strings.Contains(r.Request.Command, answer.Rock) {
					// TODO вставить фразу о запуске рока
					currentGameTracks = allTracks.Rock
					sessions[r.Session.SessionID] = userSession
					resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				} else if strings.Contains(r.Request.Command, answer.Any) {
					// TODO вставить фразу о запуске любого
					currentGameTracks = append(allTracks.NotRock, allTracks.Rock...)
					sessions[r.Session.SessionID] = userSession
					resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				} else {
					// TODO здесь надо находить жанр, похожий на названный
					resp.Text, resp.TTS = answer.IDontUnderstandYouPhrase()
				}
			} else if userSession.GameStatus == models.Playing {
				// логика во время игры
				if strings.Contains(r.Request.Command, answer.ChangeGenre) || strings.Contains(r.Request.Command, answer.AnotherGenre) {
					userSession.GameStatus = models.ChoosingGenre
					resp.Text, resp.TTS = answer.ChooseGenre, answer.ChooseGenre
				} else if userSession.MusicStarted {
					if strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Title)) && strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Artist)) {
						resp.Text, resp.TTS = answer.WinPhrase(userSession)
						userSession.MusicStarted = false
					} else if strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Title)) {
						userSession.TitleMatch = true
						if userSession.ArtistMatch {
							resp.Text, resp.TTS = answer.WinPhrase(userSession)
							userSession.MusicStarted = false
						} else {
							resp = game.WrongAnswerPlay(userSession, resp)
						}
					} else if strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Artist)) {
						userSession.ArtistMatch = true
						if userSession.TitleMatch {
							resp.Text, resp.TTS = answer.WinPhrase(userSession)
							userSession.MusicStarted = false
						} else {
							resp = game.WrongAnswerPlay(userSession, resp)
						}
					} else if userSession.NextLevelLoses {
						userSession.MusicStarted = false
						resp.Text, resp.TTS = answer.LosePhrase(userSession)
					} else {
						resp = game.WrongAnswerPlay(userSession, resp)
					}
				} else {
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
