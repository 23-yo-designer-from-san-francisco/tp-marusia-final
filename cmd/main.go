package main

import (
	"fmt"
	"guessTheSongMarusia/answer"
	"guessTheSongMarusia/game"
	"guessTheSongMarusia/models"
	"net/http"
	"strings"

	"github.com/SevereCloud/vksdk/v2/marusia"
)

// Навык "Угадай музло"
func main() {
	wh := marusia.NewWebhook()
	wh.EnableDebuging()

	mywh := marusia.NewWebhook()
	mywh.EnableDebuging()

	

	sessions := make(map[string]*models.Session)
	tracks := []models.Track{
		// {
		// 	name: "у россии три пути",
		// 	audio: map[Duration]string{
		// 		Two:  "2000512001_456239026",
		// 		Five: "2000512001_456239025",
		// 	},
		// },
		{
			Id: 1,
			Name: "Dance",
			Audio: map[models.Duration]string{
				models.Two:  "2000512001_456239030",
				models.Five: "2000512001_456239029",
			},
		},
		// {
		// 	name: "do ya think im sexy",
		// 	audio: map[Duration]string{
		// 		Two:  "2000512001_456239028",
		// 		Five: "2000512001_456239027",
		// 	},
		// },
		{	Id: 2,
			Name: "Районы кварталы",
			Audio: map[models.Duration]string{
				models.Two:  "2000512001_456239031",
				models.Five: "2000512001_456239032",
			},
		},
	}

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

			case answer.Play, answer.Playem:
				if !userSession.GameStarted {	
					sessions[r.Session.SessionID] = userSession
					resp = game.StartGame(userSession, tracks, resp)
					return resp
				}

				resp.Text, resp.TTS = answer.AlreadyPlayingPhrase()

			case answer.IDontKnow, answer.DontKnow, answer.No, answer.CantGuess, answer.ICantGuess:
				resp = game.WrongAnswerPlay(userSession, resp)

			case answer.Continue:
				if !userSession.GameStarted {
					//Ответь, что мол как продолжить, когда не начали
					resp.Text = "Ты глупый?"
					resp.TTS = "Ты глупый?"
					return
				}
				resp = game.StartGame(userSession,tracks,resp)
				return resp

			case marusia.OnInterrupt:
				resp.Text, resp.TTS = answer.GoodbyePhrase()
				resp.EndSession = true
				delete(sessions, r.Session.SessionID)

			default:
				fmt.Println("ok: ", ok, "music started: ", userSession.MusicStarted)
				fmt.Println("Command: ", r.Request.Command, "Track name: ", userSession.CurrentTrack.Name)

				if ok && userSession.MusicStarted {
					if strings.Contains(r.Request.Command, strings.ToLower(userSession.CurrentTrack.Name)) {
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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
