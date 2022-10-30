package main

import (
	"fmt"
	"guessTheSongMarusia/answer"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/marusia"
	"github.com/seehuhn/mt19937"
)

//func moreSecondsPlay()

type myPayload struct {
	Text string
	marusia.DefaultPayload
}

func getAnswerString(time Duration, audioVkId string) string {
	s := ""
	if time == 2 {
		s = "ы"
	}
	fmt.Println("getAnswerString: ", audioVkId, "Current level: ", time)
	fmt.Println(audioVkId)
	return fmt.Sprintf("Играю %d секунд%s трека. Угадаете? <speaker audio_vk_id=%s >", time, s, audioVkId)
}

func saySongInfoString(userSession *Session) string {
	return fmt.Sprintf("Это же %s", userSession.currentTrack.name)
}

func losePhrase(userSession *Session) string {
	return "Эх, вы не ^уга`дали^. Сейчас скажу вам ответ:" + 
							saySongInfoString(userSession) + "Чтобы продолжить, скажите играем"
}

func WrongAnswerPlay(userSession *Session, resp marusia.Response) marusia.Response {
	if (userSession.musicStarted && !userSession.nextLevelLoses) {
		// Тут надо инкремент, а не хардкод
		userSession.currentLevel = Five
		userSession.nextLevelLoses = true
		userSession.musicStarted = true
		resp.Text = getAnswerString(userSession.currentLevel, userSession.currentTrack.audio[userSession.currentLevel])
		return resp
	} 
		resultString := "Сейчас скажу вам ответ:" + saySongInfoString(userSession)
		resp.Text, resp.TTS = resultString, resultString
		return resp
}

type Duration int64

const (
	twoSecondsLevel = iota
	fiveSecondsLevel = iota
	tenSecondsLevel = iota
)

const (
	Two  Duration = 2
	Five          = 5
	Ten           = 10
)

type Track struct {
	id    int64
	name  string
	audio map[Duration]string
}

type Session struct {
	currentLevel   Duration
	currentPoints  int64
	currentTrack   Track
	musicStarted   bool
	nextLevelLoses bool
	playedTracks map[int]bool
}

func newSession() *Session {
	return &Session{
		currentLevel: Two,
		playedTracks: make(map[int]bool),
	}
}

func chooseTrack(userSession *Session, tracks []Track) Track {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	var randTrackID int
	for {
		rand.Seed(time.Now().Unix())
		randTrackID = rng.Int()%len(tracks)
		_, ok := userSession.playedTracks[randTrackID]
		fmt.Println(len(userSession.playedTracks))
		if !ok || len(userSession.playedTracks) == len(tracks) {
			userSession.playedTracks[randTrackID] = true
			break
		}
	}
	return tracks[randTrackID]
}

// Навык "Угадай музло"
func main() {
	wh := marusia.NewWebhook()
	wh.EnableDebuging()

	mywh := marusia.NewWebhook()
	mywh.EnableDebuging()

	

	sessions := make(map[string]*Session)
	tracks := []Track{
		// {
		// 	name: "у россии три пути",
		// 	audio: map[Duration]string{
		// 		Two:  "2000512001_456239026",
		// 		Five: "2000512001_456239025",
		// 	},
		// },
		{
			id: 1,
			name: "dance",
			audio: map[Duration]string{
				Two:  "2000512001_456239030",
				Five: "2000512001_456239029",
			},
		},
		// {
		// 	name: "do ya think im sexy",
		// 	audio: map[Duration]string{
		// 		Two:  "2000512001_456239028",
		// 		Five: "2000512001_456239027",
		// 	},
		// },
		{	id: 2,
			name: "районы кварталы",
			audio: map[Duration]string{
				Two:  "2000512001_456239031",
				Five: "2000512001_456239032",
			},
		},
	}

	// wh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
	// 	rng := rand.New(mt19937.New())
	// 	rng.Seed(time.Now().UnixNano())
	// 	userSession, ok := sessions[r.Session.SessionID]
	// 	if !ok {
	// 		userSession = Session{currentLevel: Two}
	// 		sessions[r.Session.SessionID] = userSession
	// 	}
	// 	switch r.Request.Type {
	// 	case marusia.SimpleUtterance:
	// 		switch r.Request.Command {
	// 		case marusia.OnStart:
	// 			resp.Text = "Скилл запущен"
	// 			resp.TTS = "Скилл запущен, жду команд"
	// 		case "говори":
	// 			resp.Text = "говорю"
	// 			resp.TTS = "говорю"
	// 		case "покричи":
	// 			resp.Text = "^говорю^"
	// 			resp.TTS = "^говорю^"
	// 		case "играем", "продолжить":
	// 			if !userSession.musicStarted {
	// 				rand.Seed(time.Now().Unix())
	// 				track := tracks[rng.Int()%len(tracks)]
	// 				fmt.Println("Selected track", track)
	// 				userSession.currentTrack = track
	// 				userSession.musicStarted = true
	// 				sessions[r.Session.SessionID] = userSession
	// 				resp.Text = getAnswerString(userSession.currentLevel, track.audio[userSession.currentLevel])
	// 				fmt.Println("UserSession = ", userSession)
	// 			} else {
	// 				delete(sessions, r.Session.SessionID)
	// 				resp.EndSession = true
	// 			}

	// 		case "нет", "не узнал", "не знаю":
	// 			if userSession.musicStarted && !userSession.nextLevelLoses {
	// 				// Тут надо инкремент, а не хардкод
	// 				userSession.currentLevel = Five
	// 				userSession.nextLevelLoses = true
	// 				userSession.musicStarted = true
	// 				sessions[r.Session.SessionID] = userSession
	// 				resp.Text = getAnswerString(userSession.currentLevel, userSession.currentTrack.audio[userSession.currentLevel])
	// 				return
	// 			} else {
	// 				resp.Text = "Повезет в другой раз, продолжим?"
	// 			}
	// 			resp.EndSession = true
	// 			delete(sessions, r.Session.SessionID)

	// 		case marusia.OnInterrupt:
	// 			resp.Text = "Скилл закрыт"
	// 			resp.TTS = "Пока"
	// 			resp.EndSession = true

	// 		default:
	// 			fmt.Println("ok: ", ok, "music started: ", userSession.musicStarted)
	// 			fmt.Println("Command: ", r.Request.Command, "Track name: ", userSession.currentTrack.name)
	// 			if ok && userSession.musicStarted {
	// 				if strings.ToLower(r.Request.Command) == userSession.currentTrack.name {
	// 					resp.Text = "Вы молодец, ^угадали^! <speaker audio=marusia-sounds/game-win-1>"
	// 				} else {
	// 					resp.Text = "Повезет в другой раз, ^продолжим^?"
	// 				}
	// 				resp.EndSession = true
	// 				delete(sessions, r.Session.SessionID)
	// 			} else {
	// 				resp.Text = "Неизвестная команда"
	// 				resp.TTS = "Я вас \"НЕ СОВСЕМ\" поняла"
	// 			}
	// 		}
	// 	}

	// 	return
	// })

	mywh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		userSession, ok := sessions[r.Session.SessionID]
		if !ok {
			userSession = newSession()
			sessions[r.Session.SessionID] = userSession
		}
		switch r.Request.Type {
		case marusia.SimpleUtterance:
			switch r.Request.Command {
			case marusia.OnStart:
				resp.Text = "Скилл запущен"
				resp.TTS = "Скилл запущен, жду команд"
			case answer.Playem:
				if !userSession.musicStarted {	
					userSession.currentTrack = chooseTrack(userSession, tracks)
					fmt.Println("Selected track", userSession.currentTrack)
					userSession.musicStarted = true
					sessions[r.Session.SessionID] = userSession
					answerString := getAnswerString(userSession.currentLevel, userSession.currentTrack.audio[userSession.currentLevel])
					resp.Text = answerString
					resp.TTS = answerString
					return
				}
				resp.Text = "Вы уже играете, забыли?"
				resp.TTS = "Вы уже играете, забыли?"
			case answer.IDontKnow, answer.DontKnow, answer.No, answer.CantGuess, answer.ICantGuess:
				resp = WrongAnswerPlay(userSession, resp)
				// Не завершать сессию, а предложить продолжить
				resp.EndSession = true
	 			delete(sessions, r.Session.SessionID)
			case marusia.OnInterrupt:
				resp.Text = "Скилл закрыт"
				resp.TTS = "Пока"
				resp.EndSession = true
				delete(sessions, r.Session.SessionID)
			default:
				fmt.Println("ok: ", ok, "music started: ", userSession.musicStarted)
				fmt.Println("Command: ", r.Request.Command, "Track name: ", userSession.currentTrack.name)
				if ok && userSession.musicStarted {
					if strings.Contains(r.Request.Command, userSession.currentTrack.name) {
						resp.Text = "Вы молодец! ^Уга`дали^! <speaker audio=marusia-sounds/game-win-1> Чтобы продолжить, скажите играем"
						//Тут надо очистить игру, не завершая сессию
						userSession.musicStarted = false
						userSession.currentLevel = Two
						return
					}
					if userSession.nextLevelLoses {
						resultString := losePhrase(userSession)
						resp.Text, resp.TTS = resultString, resultString
						return
					}
					resp = WrongAnswerPlay(userSession, resp)
					return
				} else {
					resp.Text = "Неизвестная команда"
					resp.TTS = "Я вас \"НЕ СОВСЕМ\" поняла"
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
