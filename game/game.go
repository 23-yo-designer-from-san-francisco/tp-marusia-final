package game

import (
	"fmt"
	"guessTheSongMarusia/answer"
	"guessTheSongMarusia/models"
	"math/rand"
	"time"

	"github.com/SevereCloud/vksdk/v2/marusia"
)

func GetAnswerString(time models.Duration, audioVkId string) (string, string) {
	var s string
	fmt.Println("Секунды", time)
	switch time {
	case 2:
		s = "две секунды"
	case 5:
		s = "пять секунд"
	case 10:
		s = "десять секунд"
	}

	fmt.Println("getAnswerString: ", audioVkId, "Current level: ", time)
	fmt.Println(audioVkId)
	resultString := fmt.Sprintf("Играю %s трека. Угадаете? <speaker audio_vk_id=%s >", s, audioVkId)
	return resultString, resultString
}

func StartGame(userSession *models.Session, tracks []models.VKTrack, resp marusia.Response, rng *rand.Rand) marusia.Response {
	userSession.CurrentTrack = ChooseTrack(userSession, tracks, rng)
	fmt.Println("Selected track", userSession.CurrentTrack)
	userSession.GameStatus = models.Playing
	userSession.MusicStarted = true
	userSession.CurrentLevel = models.Two
	userSession.NextLevelLoses = false
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
	var text, tts string
	switch userSession.CurrentLevel {
	case models.Two:
		text, tts = GetAnswerString(userSession.CurrentLevel, userSession.CurrentTrack.Duration2)
	case models.Five:
		text, tts = GetAnswerString(userSession.CurrentLevel, userSession.CurrentTrack.Duration3)
	case models.Ten:
		text, tts = GetAnswerString(userSession.CurrentLevel, userSession.CurrentTrack.Duration5)
	}
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
