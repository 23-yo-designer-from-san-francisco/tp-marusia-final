package models

import (
	"fmt"
)

//Phrase должны возвращать (string, string) -> resp.Text, resp.TTS

// Неизвестная команда
func IDontUnderstandYouPhrase() (string, string) {
	return UnknownCommand, UnknownCommandTTS
}

// Информация о загаданной песне
func SaySongInfoString(userSession *Session) string {
	return fmt.Sprintf("%s %s — %s.", ThatIs, userSession.CurrentTrack.Artist, userSession.CurrentTrack.Title)
}

func GetScoreText(userSession *Session) string {
	var score string
	if userSession.CompetitionMode {
		score = fmt.Sprintf("%s %g %s.", YourScore, userSession.CurrentPoints, "баллов")
	}
	return score
}

//Функция, меняющая state вне main. Хз, норм ли. По другому не придумал.
func CheckPlaylistFinished(userSession *Session, str string) string {
	if len(userSession.CurrentPlaylist) == 0 {
		str = fmt.Sprintf("%s %s", str, PlaylistFinished)
		userSession.GameState = PlaylistFinishedState
	}

	return str
}

// Если человек не смог угадать
func LosePhrase(userSession *Session) (string, string) {
	var str string

	userSession.Fails += 1
	str = fmt.Sprintf("%s %s %s %s %s", DontGuess, IWillSayTheAnswer,
		SaySongInfoString(userSession), GetScoreText(userSession), ToContinue)

	str = CheckPlaylistFinished(userSession, str)

	if userSession.Fails%4 == 0 && userSession.Fails != 0 {
		str = fmt.Sprintf("%s %s", str, Notify)
	}

	return str, str
}

func WinPhrase(userSession *Session) (string, string) {
	textString := fmt.Sprintf("%s %s %s %s", YouGuessText, GetScoreText(userSession), ToContinue, ToStop)
	ttsString := fmt.Sprintf("%s %s %s %s", YouGuessTTS, SaySongInfoString(userSession), ToContinue, ToStop)
	textString = CheckPlaylistFinished(userSession, textString)
	ttsString = CheckPlaylistFinished(userSession, ttsString)
	return textString, ttsString
}

// Начало Игры
func StartGamePhrase() (string, string) {
	str := fmt.Sprintf("%s %s %s %s", Hello, ToStart, ToStartCompetitive, ToStop)
	return str, str
}

func ChooseGenrePhrase() (string, string) {
	str := fmt.Sprintf("%s", ChooseGenre)
	return str, str
}

func PlayingGamePhrase() (string, string) {
	str := `Сейчас вы играете, попробуйте угадать песню, исполнителя или всё сразу. 
		Если вы хотите поменять жанр скажите "Сменить Жанр". Чтобы сдаться, скажите "Сдаюсь"`
	return str, str
}

func CompetitionPhrase() (string, string) {
	str := Competition
	return str, str
}

func AvailableGenresPhrase() (string, string) {
	str := AvailableGenres
	return str, str
}

func ChooseArtistPhrase() (string, string) {
	str := "Назовите Исполнителя, а я посмотрю знаю ли я о нём"
	return str, str
}

func CompetitionRulesPhrase() (string, string) {
	str := CompetitionRules
	return str, str
}

func PlaylistFinishedPhrase() (string, string) {
	str := PlaylistFinished
	return str, str
}

func KeyPhrasePlaylistPhrase() (string, string) {
	str := KeyPhrasePlaylist
	return str, str
}

func StandartErrorPhrase() (string, string) {
	str := ErrorHappend
	return str, str
}
