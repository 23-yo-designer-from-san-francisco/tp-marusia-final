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
	return fmt.Sprintf("%s%s — %s. ", ThatIs, userSession.CurrentTrack.Artist, userSession.CurrentTrack.Title)
}

// Если человек не смог угадать
func LosePhrase(userSession *models.Session) string {
	return DontGuess + IWillSayTheAnswer +
		SaySongInfoString(userSession) + ToContinue + ToStop
}

func WinPhrase(userSession *models.Session) string {
	return YouGuess + SaySongInfoString(userSession) + ToContinue + ToStop
}

// Начало Игры
func StartGamePhrase() (string, string) {
	str := fmt.Sprintf("%s %s", Hello, ToStart)
	return str, str
}

// Жанры
func ChooseGenrePhrase() (string, string) {
	str := fmt.Sprintf("%s %s %s", ChooseGenre, AvailableGenres, ToStop)
	return str, str
}

func ChangeGenrePhrase() (string, string) {
	str := fmt.Sprintf("%s %s", ChangeGenre, AvailableGenres)
	return str, str
}

func AlreadyPlayingPhrase() (string, string) {
	return AlreadyPlaying, AlreadyPlaying
}

func GoodbyePhrase() (string, string) {
	return GoodBye, GoodBye
}
