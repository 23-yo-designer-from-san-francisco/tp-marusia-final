package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"guessTheSongMarusia/answer"
	"guessTheSongMarusia/microservice/music/usecase"
	"guessTheSongMarusia/models"
	"math/rand"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/marusia"
)

func StartGame(userSession *models.Session, resp marusia.Response, trackCount int, mU *usecase.MusicUsecase, rng *rand.Rand) marusia.Response {
	//TODO userSession.CurrentGenre тут выбрать жанр нужный
	userSession.CurrentTrack = ChooseTrack(userSession, trackCount, mU, rng)
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

func ChooseTrack(userSession *models.Session, trackCount int, mU *usecase.MusicUsecase, rng *rand.Rand) models.VKTrack {
	var randTrackID int

	for {
		rand.Seed(time.Now().Unix())
		randTrackID = rng.Int() % trackCount
		fmt.Println("Total tracks", trackCount, " Random track: ", randTrackID)
		_, ok := userSession.PlayedTracks[randTrackID]
		fmt.Println(len(userSession.PlayedTracks))
		if !ok || len(userSession.PlayedTracks) == trackCount {
			userSession.PlayedTracks[randTrackID] = true
			break
		}
	}
	track, err := mU.GetSongById(randTrackID)
	if err != nil {
		logrus.Fatal(err)
	}
	return track
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

func SelectGenre(userSession *models.Session, command string, resp marusia.Response, trackCount int, mU *usecase.MusicUsecase,
	sessions map[string]*models.Session, allTracks models.TracksPerGenres, sessionID string, rng *rand.Rand) marusia.Response {
	if strings.Contains(command, strings.ToLower(answer.NotRock)) {
		// не рок
		// TODO вставить фразу о запуске не рока
		sessions[sessionID] = userSession
		userSession.GenreTrackCounter = 0
		userSession.CurrentGenre = answer.NotRock
		resp = StartGame(userSession, resp, trackCount, mU, rng)
	} else if strings.Contains(command, strings.ToLower(answer.Rock)) {
		// рок
		// TODO вставить фразу о запуске рока
		sessions[sessionID] = userSession
		userSession.GenreTrackCounter = 0
		userSession.CurrentGenre = answer.Rock
		resp = StartGame(userSession, resp, trackCount, mU, rng)
	} else if strings.Contains(command, strings.ToLower(answer.Any)) {
		// любой
		// TODO вставить фразу о запуске любого
		sessions[sessionID] = userSession
		userSession.GenreTrackCounter = 0
		userSession.CurrentGenre = answer.Any
		resp = StartGame(userSession, resp, trackCount, mU, rng)
	} else {
		// непонел
		// TODO здесь надо находить жанр, похожий на названный
		resp.Text, resp.TTS = answer.IDontUnderstandYouPhrase()
	}
	return resp
}
