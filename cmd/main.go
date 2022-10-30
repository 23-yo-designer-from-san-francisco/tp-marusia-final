package main

import (
	"github.com/SevereCloud/vksdk/v2/marusia"
	"net/http"
)

import "encoding/json"

type Art struct {
	URL string `json:"url"`
}

type myPayload struct {
	Text string
	marusia.DefaultPayload
}

// Навык "Угадай музло"
func main() {
	wh := marusia.NewWebhook()
	wh.EnableDebuging()

	wh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		switch r.Request.Type {
		case marusia.SimpleUtterance:
			switch r.Request.Command {
			case marusia.OnStart:
				resp.Text = "Скилл запущен"
				resp.TTS = "Скилл запущен, жду команд"
			case "картинка":
				resp.Card = marusia.NewBigImage(
					"Заголовок",
					"Описание",
					457239017,
				)
			case "картинки":
				resp.Card = marusia.NewImageList(
					457239017,
					457239018,
				)
			case "кнопки":
				resp.Text = "Держи кнопки"
				resp.TTS = "Жми на кнопки"
				resp.AddURL("ссылка", "https://vk.com")
				resp.AddButton("подсказка без нагрузки", nil)
				resp.AddButton("подсказка с нагрузкой", myPayload{
					Text: "test",
				})
			case "ссылка":
				resp.Text = marusia.CreateDeepLink(
					"e7a7d540-3928-4f11-87bf-a0de1244c096",
					map[string]string{"Text": "нагрузка из ссылки"},
				)
				resp.TTS = "Держи диплинк"
			case "пуш":
				resp.Text = `Держи пуш`
				resp.TTS = `Отправила пуш на устройство`
				resp.Push.PushText = "Hello, i am push"
			case "музон":
				resp.Text = "Музон запущен"
				resp.TTS = "Включаю Музон"
				player := marusia.AudioPlayer{
					SeekTrack: 0,
					SeekSecond: 0,
					Playlist: []marusia.AudioPlaylist{
						{
							Meta: marusia.AudioMeta{
								Title: "title",
								SubTitle: "subtitle",
								Art: Art{URL: "https://sun1-91.userapi.com/impf/wpVxuGAZV3ItESy681IpYLT9UuNt5xainEruLw/j10IeDal8cE.jpg?size=160x0&quality=90&sign=c46d43fc26c10d0623261e96cfac7f0a"},
							},
							Stream: marusia.AudioStream{
								TrackID:    "SomeText",
								SourceType: "vk",
								Source:     "-2001702405_114702405",
							},
						},
					},
				}
				resp.AudioPlayer = &player
			case marusia.OnInterrupt:
				resp.Text = "Скилл закрыт"
				resp.TTS = "Пока"
				resp.EndSession = true
			default:
				resp.Text = "Неизвестная команда"
				resp.TTS = "Я вас не поняла"
			}
		case marusia.ButtonPressed:
			var p myPayload

			err := json.Unmarshal(r.Request.Payload, &p)
			if err != nil {
				resp.Text = "Что-то пошло не так"
				return
			}

			resp.Text = "Кнопка нажата. Полезная нагрузка: " + p.Text
			resp.TTS = "Вы нажали на кнопку"
		case marusia.DeepLink:
			var p myPayload

			err := json.Unmarshal(r.Request.Payload, &p)
			if err != nil {
				resp.Text = "Что-то пошло не так"
				return
			}

			resp.Text = "Специальная ссылка. Полезная нагрузка: " + p.Text
			resp.TTS = "Вы перешли по ссылке"
		}

		return
	})

	http.HandleFunc("/", wh.HandleFunc)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
