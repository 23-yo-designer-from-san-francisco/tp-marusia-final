package game

import (
	"fmt"
	"guessTheSongMarusia/answer"
	"guessTheSongMarusia/models"
	"math/rand"
	"time"

	"github.com/SevereCloud/vksdk/v2/marusia"
)

func StartGame(userSession *models.Session, tracks []models.VKTrack, resp marusia.Response, rng *rand.Rand) marusia.Response {
	userSession.CurrentTrack = ChooseTrack(userSession, tracks, rng)
	fmt.Println("Selected track", userSession.CurrentTrack)
	userSession.GameStatus = models.Playing
	userSession.MusicStarted = true
	userSession.CurrentLevel = models.Two
	userSession.NextLevelLoses = false
	userSession.ArtistMatch = false
	userSession.TitleMatch = false
	userSession.GenreTrackCounter += 1
	resp.Text, resp.TTS = getRespTextFromLevel(userSession)
	return resp
}

func ChooseTrack(userSession *models.Session, tracks []models.VKTrack, rng *rand.Rand) models.VKTrack {
	var randTrackID int
	for {
		rand.Seed(time.Now().Unix())
		randTrackID = rng.Int() % len(tracks)
		_, ok := userSession.PlayedTracks[randTrackID]
		fmt.Println(len(userSession.PlayedTracks))
		if !ok || len(userSession.PlayedTracks) == len(tracks) {
			userSession.PlayedTracks[randTrackID] = true
			break
		}
	}
	return tracks[randTrackID]
}

func getRespTextFromLevel(userSession *models.Session) (string, string) {
	var s, text, tts, preWin, audioVkId string
	fmt.Println("Секунды", userSession.CurrentLevel)
	switch userSession.CurrentLevel {
	case models.Two:
		s = "две секунды"
		audioVkId = userSession.CurrentTrack.Duration2
	case models.Five:
		s = "следующие три секунды"
		audioVkId = userSession.CurrentTrack.Duration3
	case models.Ten:
		s = "следующие пять секунд"
		audioVkId = userSession.CurrentTrack.Duration5
	}

	if userSession.ArtistMatch {
		preWin = "Вы угадали исполнителя! А ^см`ожете^ название? "
	} else if userSession.TitleMatch {
		preWin = "Вы угадали название! А ^см`ожете^ исполнителя? "
	} else if userSession.GenreTrackCounter == 1 && userSession.CurrentLevel == models.Two {
		preWin = fmt.Sprintf("Вы выбрали жанр «%s». Чтобы поменять, скажите «сменить жанр». ", userSession.CurrentGenre)
	}

	fmt.Println("getAnswerString: ", audioVkId, "Current level: ", userSession.CurrentLevel)
	fmt.Println(audioVkId)
	text = fmt.Sprintf("%sИграю %s трека. Угадаете?", preWin, s)
	tts = fmt.Sprintf("%s <speaker audio_vk_id=%s >", text, audioVkId)
	return text, tts
}

func WrongAnswerPlay(userSession *models.Session, resp marusia.Response) marusia.Response {
	if userSession.MusicStarted && !userSession.NextLevelLoses {
		// Тут надо инкремент, а не хардкод
		if userSession.CurrentLevel == models.Two {
			userSession.CurrentLevel = models.Five
		} else if userSession.CurrentLevel == models.Five {
			userSession.CurrentLevel = models.Ten
			userSession.NextLevelLoses = true
		}
		userSession.MusicStarted = true
		resp.Text, resp.TTS = getRespTextFromLevel(userSession)
		return resp
	}
	resultString := fmt.Sprintf("%s — %s %s %s. ", answer.IWillSayTheAnswer, answer.SaySongInfoString(userSession), answer.ToContinue, answer.ToStop)
	resp.Text, resp.TTS = resultString, resultString
	return resp
}
