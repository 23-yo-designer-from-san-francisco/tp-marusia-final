package models

import (
	"fmt"
	"strings"
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
	pointsStr := "баллов"
	//От 10 до 19 оставляем баллов
	if userSession.CurrentPoints / 10 == 1 {
	} else if int(userSession.CurrentPoints) % 10 == 1 {
		pointsStr = "балл"
	} else if int(userSession.CurrentPoints) % 10 >= 2 && int(userSession.CurrentPoints) % 10 <= 4 {
		pointsStr = "балла"
	}
	if userSession.CompetitionMode {
		score = fmt.Sprintf("%s %g %s.", YourScore, userSession.CurrentPoints, pointsStr)
	}
	return score
}

func CheckPlaylistFinished(userSession *Session, str string) string {
	if len(userSession.CurrentPlaylist) == 0 {
		str = fmt.Sprintf("%s %s", str, PlaylistFinished)
		userSession.GameState = PlaylistFinishedState
		if userSession.CompetitionMode {
			str = fmt.Sprintf("%s %s %s", str, "Ключевая фраза вашего плейлиста:", strings.Title(userSession.KeyPhrase))
		}
	}

	return str
}

// Если человек не смог угадать
func LosePhrase(userSession *Session) (string, string) {
	var str string
	userSession = countPoints(userSession)
	userSession.Fails += 1
	str = fmt.Sprintf("%s %s %s %s %s", DontGuess, IWillSayTheAnswer,
		SaySongInfoString(userSession), GetScoreText(userSession), ToContinue)

	str = CheckPlaylistFinished(userSession, str)

	if userSession.Fails%4 == 0 && userSession.Fails != 0 {
		str = fmt.Sprintf("%s %s", str, Notify)
	}

	return str, str
}

func addPoints (userSession *Session, multiplier float64) (*Session) {
	switch userSession.CurrentLevel {
	case Two:
		userSession.CurrentPoints += GuessedAttempt1 * multiplier
	case Five:
		userSession.CurrentPoints += GuessedAttempt2 * multiplier
	case Ten:
		userSession.CurrentPoints += GuessedAttempt3 * multiplier
	}
	fmt.Println("Points: ", userSession.CurrentPoints)
	return userSession
}

//Должна вызываться перед winPhrase and losePhrase
func countPoints(userSession *Session) (*Session) {
	fmt.Println("ArtistMatch: ",userSession.ArtistMatch)
	fmt.Println("TitleMatch: ",userSession.TitleMatch)
	//Ничего не угадали
	if !userSession.ArtistMatch && !userSession.TitleMatch {
		return userSession
	}
	//Угадали только название (Здесь не может быть в теории artistMode)
	if userSession.TitleMatch && !userSession.ArtistMatch {
		addPoints(userSession, 0.5)
		return userSession
	}
	//Угадали только исполнителя
	if userSession.ArtistMatch && !userSession.TitleMatch {
		//Но мы в ArtistMode
		if userSession.GameMode == ArtistMode {
			return userSession
		}
		addPoints(userSession, 0.5)
		return userSession
	}
	
	//Угадано оба
	addPoints(userSession, 1)
	return userSession
}

func WinPhrase(userSession *Session) (string, string) {
	userSession = countPoints(userSession)
	fmt.Println("After Func Points: ", userSession.CurrentPoints)
	textString := fmt.Sprintf("%s %s %s %s", YouGuessText, GetScoreText(userSession), ToContinue, ToStop)
	ttsString := fmt.Sprintf("%s %s %s %s", YouGuessTTS, SaySongInfoString(userSession), ToContinue, ToStop)
	textString = CheckPlaylistFinished(userSession, textString)
	ttsString = CheckPlaylistFinished(userSession, ttsString)
	return textString, ttsString
}

// Начало Игры
func StartGamePhrase() (string, string) {
	str := fmt.Sprintf("%s %s %s %s %s", Hello, ToStart, ToStartCompetitive, ToStop, ToKeyPhrase)
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
