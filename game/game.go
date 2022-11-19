package game

import (
	"fmt"
	"guessTheSongMarusia/microservice/music/usecase"
	"guessTheSongMarusia/models"
	"math/rand"
	"time"

	"github.com/SevereCloud/vksdk/v2/marusia"
)

func StartGame(userSession *models.Session, resp marusia.Response, mU *usecase.MusicUsecase, rng *rand.Rand) marusia.Response {
	//TODO userSession.CurrentGenre тут выбрать жанр нужный
	userSession.CurrentTrack = ChooseTrack(userSession, mU, rng)
	fmt.Println("Selected track", userSession.CurrentTrack)
	fmt.Println("Selected track Artists", userSession.CurrentTrack.ArtistsWithHumanArtists)
	userSession.GameState = models.PlayingState
	userSession.MusicStarted = true
	userSession.CurrentLevel = models.Two
	userSession.NextLevelLoses = false
	userSession.ArtistMatch = false
	userSession.TitleMatch = false
	userSession.TrackCounter += 1
	resp.Text, resp.TTS = getRespTextFromLevel(userSession)
	return resp
}

func ChooseTrack(userSession *models.Session, mU *usecase.MusicUsecase, rng *rand.Rand) models.VKTrack {
	var randTrackID int
	for {
		rand.Seed(time.Now().Unix())
		tracksCount := len(userSession.CurrentPlaylist)
		randTrackID = rng.Int() % tracksCount
		fmt.Println("Total tracks", tracksCount, " Random track: ", randTrackID)
		_, ok := userSession.PlayedTracks[randTrackID]
		fmt.Println("Length of map PlayedTracks", len(userSession.PlayedTracks))
		fmt.Println("TrackCounter", userSession.TrackCounter)
		if !ok || userSession.TrackCounter >= tracksCount {
			userSession.PlayedTracks[randTrackID] = true
			break
		}
	}
	return userSession.CurrentPlaylist[randTrackID]
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

	if userSession.ArtistMatch && userSession.GameMode != models.ArtistMode {
		preWin = "Вы угадали исполнителя! А ^см`ожете^ название? "
	} else if userSession.TitleMatch {
		preWin = "Вы угадали название! А ^см`ожете^ исполнителя? "
	} else if userSession.TrackCounter == 1 && userSession.CurrentLevel == models.Two {
		if userSession.GameMode == models.ArtistMode {
			preWin = fmt.Sprintf("Вы выбрали исполнителя «%s». Вы можете в любой момент «Сменить игру», «Сменить исполнителя» или «Сменить жанр». ", userSession.CurrentGenre)
		} else {
			preWin = fmt.Sprintf("Вы выбрали жанр «%s». Вы можете в любой момент «Сменить игру», «Сменить исполнителя» или «Сменить жанр». ", userSession.CurrentGenre)
		}
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
	resultString := fmt.Sprintf("%s — %s %s. ", models.IWillSayTheAnswer, models.SaySongInfoString(userSession), models.ToContinue)
	resp.Text, resp.TTS = resultString, resultString
	return resp
}

func CloseAnswerPlay(userSession *models.Session, resp marusia.Response) marusia.Response {
	if userSession.ArtistMatch {
		if userSession.GameMode != models.ArtistMode {
			str := "Вы правильно назвали исполнителя. А сможете сказать название песни?"
			resp.Text, resp.TTS = str, str
		} else {
			str := "Исполнитель вам и так известен, а вот с названием похоже не везёт. Подумайте ещё"
			resp.Text, resp.TTS = str, str
		}

		return resp
	}

	if userSession.TitleMatch {
		str := "Вы угадали название песни. А сможете назвать исполнителя?"
		resp.Text, resp.TTS = str, str
	}
	return resp
}

func SelectGenre(userSession *models.Session, command string, resp marusia.Response, mU *usecase.MusicUsecase,
	sessionID string, rng *rand.Rand) marusia.Response {
	var tracks []models.VKTrack
	var err error
	var genre string
	if command == "любой" {
		tracks, err = mU.GetAllMusic()
		genre = "Любой"
	} else {
		tracks, _ = mU.GetMusicByGenre(command)
		genre, err = mU.GetGenreFromHumanGenre(command)
	}

	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil || len(tracks) == 0 {
		str := "Извините, я не нашла нужный жанр, либо просто вас не поняла. Попробуйте ещё"
		resp.Text, resp.TTS = str, str
		return resp
	}

	userSession.TrackCounter = 0
	userSession.CurrentGenre = genre
	userSession.GameMode = models.GenreMode
	userSession.CurrentPlaylist = tracks
	resp = StartGame(userSession, resp, mU, rng)
	return resp
}

func SelectArtist(userSession *models.Session, command string, resp marusia.Response, mU *usecase.MusicUsecase,
		sessionID string, rng *rand.Rand) marusia.Response {
	tracks, artist, err := mU.GetSongsByArtist(command)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil || len(tracks) == 0 {
		str := "Извините, я не знаю о таком исполнителя. Назовите кого-нибудь ещё."
		resp.Text, resp.TTS = str, str
		return resp
	}

	userSession.TrackCounter = 0
	userSession.GameMode = models.ArtistMode
	userSession.CurrentGenre = artist
	userSession.CurrentPlaylist = tracks
	resp = StartGame(userSession, resp, mU, rng)
	return resp
}
