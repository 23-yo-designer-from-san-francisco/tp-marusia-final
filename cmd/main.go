package main

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/marusia"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

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

type Duration int64

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
}

// Навык "Угадай музло"
func main() {
	wh := marusia.NewWebhook()
	wh.EnableDebuging()

	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	sessions := make(map[string]Session)
	tracks := []Track{
		// {
		// 	name: "у россии три пути",
		// 	audio: map[Duration]string{
		// 		Two:  "2000512001_456239026",
		// 		Five: "2000512001_456239025",
		// 	},
		// },
		{
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
		{
			name: "районы кварталы",
			audio: map[Duration]string{
				Two:  "2000512001_456239031",
				Five: "2000512001_456239032",
			},
		},
	}

	wh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		userSession, ok := sessions[r.Session.SessionID]
		if !ok {
			userSession = Session{currentLevel: Two}
			sessions[r.Session.SessionID] = userSession
		}
		switch r.Request.Type {
		case marusia.SimpleUtterance:
			switch r.Request.Command {
			case marusia.OnStart:
				resp.Text = "Скилл запущен"
				resp.TTS = "Скилл запущен, жду команд"
			case "играем", "продолжить":
				if !userSession.musicStarted {
					rand.Seed(time.Now().Unix())
					track := tracks[rng.Int()%len(tracks)]
					fmt.Println("Selected track", track)
					userSession.currentTrack = track
					userSession.musicStarted = true
					sessions[r.Session.SessionID] = userSession
					resp.Text = getAnswerString(userSession.currentLevel, track.audio[userSession.currentLevel])
				} else {
					delete(sessions, r.Session.SessionID)
					resp.EndSession = true
				}
			case "нет", "не узнал", "не знаю":
				if userSession.musicStarted && !userSession.nextLevelLoses {
					// Тут надо инкремент, а не хардкод
					userSession.currentLevel = Five
					userSession.nextLevelLoses = true
					userSession.musicStarted = true
					sessions[r.Session.SessionID] = userSession
					resp.Text = getAnswerString(userSession.currentLevel, userSession.currentTrack.audio[userSession.currentLevel])
				} else {
					resp.Text = "Повезет в другой раз."
				}
				resp.EndSession = true
				delete(sessions, r.Session.SessionID)
			case marusia.OnInterrupt:
				resp.Text = "Скилл закрыт"
				resp.TTS = "Пока"
				resp.EndSession = true
			default:
				fmt.Println("ok: ", ok, "music started: ", userSession.musicStarted)
				fmt.Println("Command: ", r.Request.Command, "Track name: ", userSession.currentTrack.name)
				if ok && userSession.musicStarted {
					if strings.ToLower(r.Request.Command) == userSession.currentTrack.name {
						resp.Text = "Вы молодец, угадали!"
					} else {
						resp.Text = "Повезет в другой раз."
					}
					resp.EndSession = true
					delete(sessions, r.Session.SessionID)
				} else {
					resp.Text = "Неизвестная команда"
					resp.TTS = "Я вас \"НЕ СОВСЕМ\" поняла"
				}
			}
		}

		return
	})

	http.HandleFunc("/", wh.HandleFunc)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
