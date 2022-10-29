package main

import (
	"github.com/SevereCloud/vksdk/v2/marusia"
	"net/http"
)

type myPayload struct {
	Text string
	marusia.DefaultPayload
}

func getTextStringFromId(audioVkId string) string {
	return "Что сейчас играет? <speaker audio_vk_id=" + audioVkId + ">"
}

type Duration int64

const (
	Two Duration = iota
	Five
	Ten
)

type Session struct {
	currentLevel     Duration
	currentPoints    int64
	currentTrackName string
}

// Навык "Угадай музло"
func main() {
	wh := marusia.NewWebhook()
	wh.EnableDebuging()

	sessions := make(map[string]Session)

	wh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		userSession, ok := sessions[r.Session.SessionID]
		if !ok {
			sessions[r.Session.SessionID] = Session{currentLevel: Two, currentPoints: 0, currentTrackName: ""}
		}
		switch r.Request.Type {
		case marusia.SimpleUtterance:
			switch r.Request.Command {
			case marusia.OnStart:
				resp.Text = "Скилл запущен"
				resp.TTS = "Скилл запущен, жду команд"
			case "+":
				userSession.currentTrackName = "gspd россия"
				sessions[r.Session.SessionID] = userSession
				resp.Text = getTextStringFromId("2000512001_456239026")
			case marusia.OnInterrupt:
				resp.Text = "Скилл закрыт"
				resp.TTS = "Пока"
				resp.EndSession = true
			default:
				if ok {
					if r.Request.Command == userSession.currentTrackName {
						resp.Text = "Вы молодец, угадали!"
					} else {
						resp.Text = "Повезет в другой раз."
					}
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
