package answer

import (
	"fmt"
	"guessTheSongMarusia/models"
)

//Phrase должны возвращать (string, string) -> resp.Text, resp.TTS

// Неизвестная команда
func IDontUnderstandYouPhrase() (string, string) {
	return UnknownCommand, UnknownCommandTTS
}

// Информация о загаданной песне
func SaySongInfoString(userSession *models.Session) string {
	return fmt.Sprintf("%s %s — %s.", ThatIs, userSession.CurrentTrack.Artist, userSession.CurrentTrack.Title)
}

// Если человек не смог угадать
func LosePhrase(userSession *models.Session) (string, string) {
	var str string
	userSession.Fails += 1
	if userSession.Fails >= 3 {
		str = fmt.Sprintf("%s %s %s %s %s", DontGuess, IWillSayTheAnswer,
			SaySongInfoString(userSession), ToContinue, GenreNotify)
	} else {
		str = fmt.Sprintf("%s %s %s %s %s", DontGuess, IWillSayTheAnswer,
			SaySongInfoString(userSession), ToContinue, ToStop)
	}

	return str, str
}

func WinPhrase(userSession *models.Session) (string, string) {
	textString := fmt.Sprintf("%s %s %s", YouGuessText, ToContinue, ToStop)
	ttsString := fmt.Sprintf("%s %s %s %s", YouGuessTTS, SaySongInfoString(userSession), ToContinue, ToStop)
	return textString, ttsString
}

// Начало Игры
func StartGamePhrase() (string, string) {
	str := fmt.Sprintf("%s %s %s %s", Hello, ToStart, ToStartCompetitive, ToStop)
	return str, str
}
