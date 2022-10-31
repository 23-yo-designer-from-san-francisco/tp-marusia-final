package main

import (
	"encoding/json"
	"fmt"
	"guessTheSongMarusia/answer"
	"guessTheSongMarusia/game"
	"guessTheSongMarusia/models"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SevereCloud/vksdk/v2/marusia"
)

//const (
//	jsonTracks = "[{\"title\":\"Psycho\",\"artist\":\"Post Malone, Ty Dolla $ign\",\"audioVkId\":\"2000512001_456239049\",\"duration\":10},{\"title\":\"Sunflower\",\"artist\":\"Post Malone, Swae Lee\",\"audioVkId\":\"2000512001_456239048\",\"duration\":10},{\"title\":\"Congratulations\",\"artist\":\"Post Malone, Quavo\",\"audioVkId\":\"2000512001_456239047\",\"duration\":10},{\"title\":\"Take What You Want\",\"artist\":\"Post Malone, Ozzy Osbourne, Travis Scott\",\"audioVkId\":\"2000512001_456239045\",\"duration\":10},{\"title\":\"rockstar\",\"artist\":\"Post Malone, 21 Savage\",\"audioVkId\":\"2000512001_456239044\",\"duration\":10},{\"title\":\"Wow.\",\"artist\":\"Post Malone\",\"audioVkId\":\"2000512001_456239043\",\"duration\":10},{\"title\":\"White Iverson\",\"artist\":\"Post Malone\",\"audioVkId\":\"2000512001_456239041\",\"duration\":10},{\"title\":\"Rich &amp; Sad\",\"artist\":\"Post Malone\",\"audioVkId\":\"2000512001_456239040\",\"duration\":10},{\"title\":\"Motley Crew\",\"artist\":\"Post Malone\",\"audioVkId\":\"2000512001_456239039\",\"duration\":10},{\"title\":\"Go Flex\",\"artist\":\"Post Malone\",\"audioVkId\":\"2000512001_456239037\",\"duration\":10},{\"title\":\"Circles\",\"artist\":\"Post Malone\",\"audioVkId\":\"2000512001_456239036\",\"duration\":10},{\"title\":\"Better Now\",\"artist\":\"Post Malone\",\"audioVkId\":\"2000512001_456239035\",\"duration\":10},{\"title\":\"Life&#39;s A Mess II\",\"artist\":\"Juice WRLD, Clever, Post Malone\",\"audioVkId\":\"2000512001_456239034\",\"duration\":10},{\"title\":\"Wolves\",\"artist\":\"Big Sean, Post Malone\",\"audioVkId\":\"2000512001_456239033\",\"duration\":10},{\"title\":\"Районы-Кварталы(5с)\",\"artist\":\"Звери\",\"audioVkId\":\"2000512001_456239032\",\"duration\":5},{\"title\":\"Районы-Кварталы(2с)\",\"artist\":\"Звери\",\"audioVkId\":\"2000512001_456239031\",\"duration\":2},{\"title\":\"D.A.N.C.E (2s)\",\"artist\":\"Justice\",\"audioVkId\":\"2000512001_456239030\",\"duration\":2},{\"title\":\"D.A.N.C.E (5s)\",\"artist\":\"Justice\",\"audioVkId\":\"2000512001_456239029\",\"duration\":5},{\"title\":\"Do Ya Think I&#39;m Sexy (2s)\",\"artist\":\"Rod Stewart\",\"audioVkId\":\"2000512001_456239028\",\"duration\":2},{\"title\":\"Do Ya Think I&#39;m Sexy (5s)\",\"artist\":\"Rod Stewart\",\"audioVkId\":\"2000512001_456239027\",\"duration\":5}]"
//)

// Навык "Угадай музло"
func main() {
	wh := marusia.NewWebhook()
	wh.EnableDebuging()

	mywh := marusia.NewWebhook()
	mywh.EnableDebuging()

	b, err := os.ReadFile("/Users/l.belyaev/marusia/cmd/music.json")
	if err != nil {
		fmt.Print(err)
	}
	jsonTracks := string(b)
	var tracks []models.VKTrack
	if err := json.Unmarshal([]byte(jsonTracks), &tracks); err != nil {
		log.Fatalln(err.Error())
	}

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
				resp = game.StartGame(userSession, tracks, resp)
				return resp

			case marusia.OnInterrupt:
				resp.Text, resp.TTS = answer.GoodbyePhrase()
				resp.EndSession = true
				delete(sessions, r.Session.SessionID)

			case "тест":
				var Tracks []models.VKTrack
				err := json.Unmarshal([]byte(jsonTracks), &Tracks)
				resp.Text = "test"
				resp.TTS = "test"
				if err != nil {
					fmt.Println(err.Error())
					return resp
				}
				fmt.Println(Tracks)
				return resp

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
