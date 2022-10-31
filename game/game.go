package game

import (
	"fmt"
	"guessTheSongMarusia/answer"
	"guessTheSongMarusia/models"
	"math/rand"
	"time"

	"github.com/SevereCloud/vksdk/v2/marusia"
	"github.com/seehuhn/mt19937"
)

func GetAnswerString(time models.Duration, audioVkId string) (string, string) {
	s := ""
	if time == 2 {
		s = "ы"
	}
	fmt.Println("getAnswerString: ", audioVkId, "Current level: ", time)
	fmt.Println(audioVkId)
	resultString := fmt.Sprintf("Играю %d секунд%s трека. Угадаете? <speaker audio_vk_id=%s >", time, s, audioVkId)
	return resultString, resultString
}

func StartGame(userSession *models.Session, tracks []models.Track, resp marusia.Response) (marusia.Response) {
	userSession.CurrentTrack = ChooseTrack(userSession, tracks)
	fmt.Println("Selected track", userSession.CurrentTrack)
	userSession.GameStarted = true
	userSession.MusicStarted = true
	userSession.CurrentLevel = models.Two
	userSession.NextLevelLoses = false
	resp.Text, resp.TTS = GetAnswerString(userSession.CurrentLevel, userSession.CurrentTrack.Audio[userSession.CurrentLevel])
	return resp
}

func ChooseTrack(userSession *models.Session, tracks []models.Track) models.Track {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	var randTrackID int
	for {
		rand.Seed(time.Now().Unix())
		randTrackID = rng.Int()%len(tracks)
		_, ok := userSession.PlayedTracks[randTrackID]
		fmt.Println(len(userSession.PlayedTracks))
		if !ok || len(userSession.PlayedTracks) == len(tracks) {
			userSession.PlayedTracks[randTrackID] = true
			break
		}
	}
	return tracks[randTrackID]
}

func WrongAnswerPlay(userSession *models.Session, resp marusia.Response) marusia.Response {
	if (userSession.MusicStarted && !userSession.NextLevelLoses) {
		// Тут надо инкремент, а не хардкод
		userSession.CurrentLevel = models.Five
		userSession.NextLevelLoses = true
		userSession.MusicStarted = true
		resp.Text, resp.TTS = GetAnswerString(userSession.CurrentLevel, userSession.CurrentTrack.Audio[userSession.CurrentLevel])
		return resp
	} 
		resultString := answer.IWillSayTheAnswer  + answer.SaySongInfoString(userSession) + answer.ToContinue + answer.ToStop
		resp.Text, resp.TTS = resultString, resultString
		return resp
}