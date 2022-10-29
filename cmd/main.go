package main

import (
	"github.com/SevereCloud/vksdk/v2/marusia"
	"net/http"
)

import "encoding/json"

type myPayload struct {
	Text string
	marusia.DefaultPayload
}

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

	err := http.ListenAndServe(":43225", nil)
	if err != nil {
		return
	}
}
