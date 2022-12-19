package models

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/seehuhn/mt19937"
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
	if userSession.CurrentPoints/10 == 1 {
	} else if userSession.CurrentPoints%10 == 1 {
		pointsStr = "балл"
	} else if userSession.CurrentPoints%10 >= 2 && userSession.CurrentPoints%10 <= 4 {
		pointsStr = "балла"
	}
	if userSession.CompetitionMode {
		score = fmt.Sprintf("%s %d %s. ", YourScore, userSession.CurrentPoints, pointsStr)
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
	} else {
		str += ToContinue + " "
		str += ToStop
	}

	return str
}

// Если человек не смог угадать
func LosePhrase(userSession *Session) (string, string) {
	var str string
	userSession = countPoints(userSession)
	userSession.Fails += 1
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	didntGuessPhrase := YouDidntGuessTexts[rng.Int63()%int64(len(YouDidntGuessTexts))]
	str = fmt.Sprintf("%s %s %s", didntGuessPhrase,
		SaySongInfoString(userSession), GetScoreText(userSession))

	ttsString := fmt.Sprintf("%s %s <speaker audio_vk_id=%s > %s", didntGuessPhrase,
		SaySongInfoString(userSession), userSession.CurrentTrack.Duration5, GetScoreText(userSession))

	str = CheckPlaylistFinished(userSession, str)
	ttsString = CheckPlaylistFinished(userSession, ttsString)
	if userSession.Fails%4 == 0 && userSession.Fails != 0 {
		str = fmt.Sprintf("%s %s", str, Notify)
	}

	return str, ttsString
}

func addPoints(userSession *Session, divider int) *Session {
	switch userSession.CurrentLevel {
	case Three:
		userSession.CurrentPoints += GuessedAttempt1 / divider
	case Five:
		userSession.CurrentPoints += GuessedAttempt2 / divider
	case Eight:
		userSession.CurrentPoints += GuessedAttempt3 / divider
	}
	fmt.Println("Points: ", userSession.CurrentPoints)
	return userSession
}

// Должна вызываться перед winPhrase and losePhrase
func countPoints(userSession *Session) *Session {
	fmt.Println("ArtistMatch: ", userSession.ArtistMatch)
	fmt.Println("TitleMatch: ", userSession.TitleMatch)
	//Ничего не угадали
	if !userSession.ArtistMatch && !userSession.TitleMatch {
		return userSession
	}
	//Угадали только название (Здесь не может быть в теории artistMode)
	if userSession.TitleMatch && !userSession.ArtistMatch {
		addPoints(userSession, 2)
		return userSession
	}
	//Угадали только исполнителя
	if userSession.ArtistMatch && !userSession.TitleMatch {
		//Но мы в ArtistMode
		if userSession.GameMode == ArtistMode {
			return userSession
		}
		addPoints(userSession, 2)
		return userSession
	}

	//Угадано оба
	addPoints(userSession, 1)
	return userSession
}

func WinPhrase(userSession *Session) (string, string) {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	guessedPhrase := YouGuessedTexts[rng.Int63()%int64(len(YouGuessedTexts))]
	userSession = countPoints(userSession)
	fmt.Println("After Func Points: ", userSession.CurrentPoints)
	textString := fmt.Sprintf("%s %s", guessedPhrase, GetScoreText(userSession))
	ttsString := textString
	if !userSession.CompetitionMode {
		textString = fmt.Sprintf("%sЭто %s — %s. ", textString, userSession.CurrentTrack.Artist, userSession.CurrentTrack.Title)
	}
	textString = CheckPlaylistFinished(userSession, textString)
	ttsString = CheckPlaylistFinished(userSession, ttsString)
	listenPhrase := LetsListenTrack[rng.Int63()%int64(len(LetsListenTrack))]
	ttsString = fmt.Sprintf("%s %s <speaker audio_vk_id=%s >", ttsString, listenPhrase, userSession.CurrentTrack.Duration8)
	return textString, ttsString
}

// Начало Игры
func StartGamePhrase() (string, string) {
	str := fmt.Sprintf("%s %s %s %s", ToStart, ToStartCompetitive, ToStop, ToKeyPhrase)
	return str, str
}

func ChooseGenrePhrase() (string, string) {
	str := fmt.Sprintf("%s", ChooseGenre)
	return str, ChooseGenreTTS
}

func PlayingGamePhrase() (string, string) {
	str := `Сейчас вы играете, попробуйте угадать песню, исполнителя или всё сразу. 
		Если вы хотите поменять жанр, скажите "Сменить Жанр". Чтобы сдаться, скажите "Сдаюсь".`
	return str, str
}

func CompetitionPhrase() (string, string) {
	str := "Вы выбрали «Соревнование». Хотите прочитаю правила или сразу «Играем»?"
	return str, str
}

func AvailableGenresPhrase() (string, string) {
	str := AvailableGenres
	return str, str
}

func ChooseArtistPhrase() (string, string) {
	str := "Назовите Исполнителя, а я посмотрю, знаю ли я о нём."
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
