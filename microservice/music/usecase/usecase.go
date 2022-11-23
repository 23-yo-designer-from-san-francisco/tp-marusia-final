package usecase

import (
	"guessTheSongMarusia/microservice/music"
	"guessTheSongMarusia/models"
	log "guessTheSongMarusia/pkg/logger"
	"html"
	"regexp"
	"strings"
)

type MusicUsecase struct {
	musicRepository music.Repository
}

func NewMusicUsecase(musicR music.Repository) *MusicUsecase {
	return &MusicUsecase{
		musicRepository: musicR,
	}
}

func removeCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}

	return strings.Map(filter, input)
}

func basicValidation(name string) string {
	name = removeCharacters(name, ".,!?")
	//Удаляет скобки и содержимое
	reg := regexp.MustCompile(`\([^)]*\)`)
	name = reg.ReplaceAllString(name, "")
	name = strings.Replace(name, "$", "s", -1)
	name = strings.Replace(name, "é", "e", -1)
	name = strings.Replace(name, "-", " ", -1)
	return strings.ToLower(name)
}

func ampersandReplacing(name string) []string {
	var resultNames []string
	if strings.Contains(name, "&") {
		name = strings.ToLower(name)
		resultNames = append(resultNames, strings.Replace(name, "&", "and", -1))
		resultNames = append(resultNames, strings.Replace(name, "&", "и", -1))
		resultNames = append(resultNames, strings.Replace(name, " &", "", -1))
	}
	return resultNames
}

func validateMusicTrackTitle(title string) []string {
	var resultTracks []string
	resultTracks = append(resultTracks, basicValidation(title))
	resultTracks = append(resultTracks, ampersandReplacing(title)...)
	return resultTracks
}

func validateMusicArtist(artist string) []string {
	var resultArtists []string
	resultArtists = append(resultArtists, basicValidation(artist))
	resultArtists = append(resultArtists, ampersandReplacing(artist)...)
	return resultArtists
}

func getArtists(artist string, currentMap map[string][]string) map[string][]string {
	var artistsMap map[string][]string
	if currentMap == nil {
		artistsMap = make(map[string][]string)
	} else {
		artistsMap = currentMap
	}
	//VK basic format as i see
	artists := strings.Split(artist, "feat.")
	for _, artistSplit := range artists {
		separator := ","
		if strings.Contains(artistSplit, "&") {
			separator = "&"
		}
		resultArtists := strings.Split(artistSplit, separator)
		for _, resultArtist := range resultArtists {
			//На всякий от пробелов
			resultArtist = strings.TrimSpace(resultArtist)
			humanArtists := validateMusicArtist(resultArtist)
			if currentValue, ok := artistsMap[resultArtist]; ok {
				artistsMap[resultArtist] = append(currentValue, humanArtists...)
			} else {
				artistsMap[resultArtist] = humanArtists
			}
		}
	}
	return artistsMap
}

func validateMusicTrack(track *models.VKTrack) (*models.VKTrack, bool) {
	track.Title = html.UnescapeString(track.Title)
	track.HumanTitles = append(track.HumanTitles, validateMusicTrackTitle(track.Title)...)

	track.Artist = html.UnescapeString(track.Artist)
	track.ArtistsWithHumanArtists = getArtists(track.Artist, track.ArtistsWithHumanArtists)
	return track, true
}

func (mU *MusicUsecase) CreateAllMusic(Tracks []models.VKTrack) error {
	for _, track := range Tracks {
		track, ok := validateMusicTrack(&track)
		if !ok {
			continue
		}
		log.Debug(track.Title)
		log.Debug(track.Artist)
		err := mU.musicRepository.CreateTrack(track)
		if err != nil {
			log.Error()
		}
	}
	return nil
}

func (mU *MusicUsecase) GetSongsByArtist(artist string) ([]models.VKTrack, string, error) {
	return mU.musicRepository.GetSongsByArtist(artist)
}

func (mU *MusicUsecase) GetGenres() ([]string, error) {
	return mU.musicRepository.GetGenres()
}

func (mU *MusicUsecase) GetMusicByGenre(genre string) ([]models.VKTrack, error) {
	return mU.musicRepository.GetMusicByGenre(genre)
}

func (mU *MusicUsecase) GetAllMusic() ([]models.VKTrack, error) {
	return mU.musicRepository.GetAllMusic()
}

func (mU *MusicUsecase) GetGenreFromHumanGenre(humanGenre string) (string, error) {
	return mU.musicRepository.GetGenreFromHumanGenre(humanGenre)
}

func (mU *MusicUsecase) GetRandomMusic(limit int) ([]models.VKTrack, error) {
	return mU.musicRepository.GetRandomMusic(limit)
} 

func (mU *MusicUsecase) GetRandomMusicByGenre(limit int, humanGenre string) ([]models.VKTrack, error) {
	return mU.musicRepository.GetRandomMusicByGenre(limit, humanGenre)

}

func (mU *MusicUsecase) GetRandomMusicByArtist(limit int, humanArtist string) ([]models.VKTrack, string, error) {
	return mU.musicRepository.GetRandomMusicByArtist(limit, humanArtist)
}