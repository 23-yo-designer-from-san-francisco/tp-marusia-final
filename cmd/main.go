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
		}

		switch r.Request.Type {
		case marusia.SimpleUtterance:
			switch r.Request.Command {

			case marusia.OnStart, answer.Greeting:
				resp.Text, resp.TTS = answer.StartGamePhrase()

			case answer.Play, answer.Playem, answer.Begin:
				resp.Text, resp.TTS = answer.ChooseGenrePhrase()
				return resp

			case answer.Rock:
				currentGameTracks = allTracks.Rock
				sessions[r.Session.SessionID] = userSession
				resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				return resp

			case answer.Rap:
				currentGameTracks = allTracks.Rock
				sessions[r.Session.SessionID] = userSession
				resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				return resp

			case answer.Any:
				currentGameTracks = append(allTracks.Rap, allTracks.Rock...)
				sessions[r.Session.SessionID] = userSession
				resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				return resp

			case answer.IDontKnow, answer.DontKnow, answer.No, answer.CantGuess, answer.ICantGuess:
				resp = game.WrongAnswerPlay(userSession, resp)

			case answer.Continue:
				resp = game.StartGame(userSession, currentGameTracks, resp, rng)
				return resp

			case marusia.OnInterrupt:
				if userSession.GameStarted {
					resp.Text, resp.TTS = answer.ChangeGenrePhrase()
					return resp
				}
				resp.Text, resp.TTS = answer.GoodbyePhrase()
				resp.EndSession = true
				delete(sessions, r.Session.SessionID)

			default:
				fmt.Println("ok: ", ok, "music started: ", userSession.MusicStarted)
				fmt.Println("Command: ", r.Request.Command, "Track name: ", userSession.CurrentTrack.Title)

				if ok && userSession.MusicStarted {
					if strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Title)) {
						resp.Text = answer.WinPhrase(userSession)
						return
					}

					if userSession.NextLevelLoses {
						resultString := answer.LosePhrase(userSession)
						resp.Text, resp.TTS = resultString, resultString
						return
					}

					resp = game.WrongAnswerPlay(userSession, resp)
					return

				} else {
					resp.Text, resp.TTS = answer.IDontUnderstandYouPhrase()
					return
				}
			}
		}
		return
	})

	http.HandleFunc("/", mywh.HandleFunc)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
