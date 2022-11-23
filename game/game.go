package game

import (
	"fmt"
	"guessTheSongMarusia/microservice/music"
	"guessTheSongMarusia/microservice/playlist"
	"guessTheSongMarusia/models"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/SevereCloud/vksdk/v2/marusia"
)

const TRACKS_IN_RAND_PLAYLIST = 8

func StartGame(userSession *models.Session, resp marusia.Response) marusia.Response {
	//TODO userSession.CurrentGenre тут выбрать жанр нужный
	userSession.CurrentTrack = NextTrack(userSession)
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

func NextTrack(userSession *models.Session) models.VKTrack {
	last := len(userSession.CurrentPlaylist) - 1
	logrus.Info("Last index is ", last)
	track := userSession.CurrentPlaylist[last]
	userSession.CurrentPlaylist = userSession.CurrentPlaylist[:last]
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
		if userSession.CompetitionMode {
			preWin = fmt.Sprintf("%s «%s». %s", "Ключевая фраза вашего плейлиста:", strings.Title(userSession.KeyPhrase), preWin)
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

func GenerateRandomPlaylist(userSession *models.Session, resp marusia.Response, pU playlist.Usecase, mU music.Usecase,
	nouns []string, adjectives []string, rng *rand.Rand) marusia.Response {
	tracks, err := mU.GetAllMusic()
	if err != nil {
		logrus.Error(err)
	}
	rng.Seed(time.Now().UnixNano())
	rng.Shuffle(len(tracks), func(i, j int) {
		tracks[i], tracks[j] = tracks[j], tracks[i]
	})
	userSession.TrackCounter = 0
	userSession.CurrentGenre = "Любой"
	userSession.GameMode = models.GenreMode
	userSession.CurrentPlaylist = tracks
	playlistTracks := tracks[:TRACKS_IN_RAND_PLAYLIST]
	logrus.Debug("Saving playlist")
	for {
		err = pU.SavePlaylist(GeneratePlaylistName(nouns, adjectives, rng), playlistTracks)
		logrus.Debug("REPEAT")
		if err == nil {
			break
		}
	}

	resp = StartGame(userSession, resp)
	return resp
}

func SelectGenre(userSession *models.Session, command string, nouns []string, adjectives []string,
	resp marusia.Response, mU music.Usecase, pU playlist.Usecase, rng *rand.Rand) marusia.Response {
	var tracks []models.VKTrack
	var err error
	var genre string
	if command == "любой" {
		if userSession.CompetitionMode {
			tracks, err = mU.GetRandomMusic(TRACKS_IN_RAND_PLAYLIST)
			if err != nil {
				resp.Text, resp.TTS = userSession.GameState.SayErrorPhrase()
				return resp
			}
			keyPhrase := GeneratePlaylistName(nouns, adjectives, rng)
			logrus.Debug("Saving playlist in SelectGenre (any)")
			for {
				err = pU.SavePlaylist(GeneratePlaylistName(nouns, adjectives, rng), tracks)
				logrus.Debug("REPEAT")
				if err == nil {
					break
				}
			}
			if err != nil {
				resp.Text, resp.TTS = userSession.GameState.SayErrorPhrase()
				return resp
			}
			userSession.KeyPhrase = keyPhrase
		} else {
			tracks, err = mU.GetAllMusic()
		}
		genre = "Любой"
	} else {
		if userSession.CompetitionMode {
			tracks, _ = mU.GetRandomMusicByGenre(TRACKS_IN_RAND_PLAYLIST, command)

			logrus.Debug("Saving playlist in SelectGenre")
			var keyPhrase string
			for {
				keyPhrase = GeneratePlaylistName(nouns, adjectives, rng)
				err = pU.SavePlaylist(GeneratePlaylistName(nouns, adjectives, rng), tracks)
				logrus.Debug("REPEAT")
				if err == nil {
					break
				}
			}
			userSession.KeyPhrase = keyPhrase
		} else {
			tracks, _ = mU.GetMusicByGenre(command)
		}
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
	rng.Seed(time.Now().UnixNano())
	rng.Shuffle(len(userSession.CurrentPlaylist), func(i, j int) {
		userSession.CurrentPlaylist[i], userSession.CurrentPlaylist[j] = userSession.CurrentPlaylist[j], userSession.CurrentPlaylist[i]
	})
	resp = StartGame(userSession, resp)
	return resp
}

func SelectArtist(userSession *models.Session, command string, nouns []string, adjectives []string,
	resp marusia.Response, mU music.Usecase, pU playlist.Usecase, rng *rand.Rand) marusia.Response {
	var tracks []models.VKTrack
	var artist string
	var err error
	if userSession.CompetitionMode {
		tracks, artist, err = mU.GetRandomMusicByArtist(TRACKS_IN_RAND_PLAYLIST, command)
		if err != nil {
			resp.Text, resp.TTS = userSession.GameState.SayErrorPhrase()
			return resp
		}

		logrus.Debug("Saving playlist in SelectArtist")
		var keyPhrase string
		for {
			keyPhrase := GeneratePlaylistName(nouns, adjectives, rng)
			err = pU.SavePlaylist(keyPhrase, tracks)
			logrus.Debug("REPEAT")
			if err == nil {
				break
			}
		}
		userSession.KeyPhrase = keyPhrase
		if err != nil {
			resp.Text, resp.TTS = userSession.GameState.SayErrorPhrase()
			return resp
		}
	} else {
		tracks, artist, err = mU.GetSongsByArtist(command)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil || len(tracks) == 0 {
		str := "Извините, я не знаю о таком исполнителя. Назовите кого-нибудь ещё."
		resp.Text, resp.TTS = str, str
		return resp
	}

	userSession.TrackCounter = 0
	userSession.CurrentGenre = artist
	userSession.GameMode = models.ArtistMode
	userSession.CurrentPlaylist = tracks
	rng.Seed(time.Now().UnixNano())
	rng.Shuffle(len(userSession.CurrentPlaylist), func(i, j int) {
		userSession.CurrentPlaylist[i], userSession.CurrentPlaylist[j] = userSession.CurrentPlaylist[j], userSession.CurrentPlaylist[i]
	})
	resp = StartGame(userSession, resp)
	return resp
}
