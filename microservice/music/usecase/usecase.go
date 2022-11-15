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
	name = strings.Replace(name, "-"," ", -1)
	//Если же что-то с чём-то, то нам такое не нужно
	return strings.ToLower(name)
}

func validateMusicTrackTitle(title string) ([]string) {
	var resultTracks []string
	//Можно добавлять свои валидации(Что удобно по идее) и доставать title
	title = basicValidation(title)

	resultTracks = append(resultTracks, title)
	return resultTracks
}

func validateMusicArtist(artist string) ([]string) {
	var resultArtists []string
	artist = basicValidation(artist)

	resultArtists = append(resultArtists, artist)
	return resultArtists
}

func getArtists(artist string) map[string][]string {
	artistsMap := make(map[string][]string)
	//VK basic format as i see
	artists := strings.Split(artist, "feat.")
	for _, artistSplit := range artists {
		resultArtists := strings.Split(artistSplit, ",")
		for _, resultArtist := range resultArtists {
			//На всякий от пробелов
			if resultArtist[0] == ' ' {
				resultArtist = resultArtist[1:]
			}
			if resultArtist[len(resultArtist)-1] == ' ' {
				resultArtist = resultArtist[:len(resultArtist)-1]
			}

			humanArtists :=  validateMusicArtist(resultArtist)
			artistsMap[resultArtist] = humanArtists
		}
	}
	return artistsMap
}



func validateMusicTrack(track *models.VKTrack) (*models.VKTrack, bool) {
	track.Title = html.UnescapeString(track.Title)
	track.HumanTitles = validateMusicTrackTitle(track.Title)

	track.Artist = html.UnescapeString(track.Artist)
	track.ArtistsWithHumanArtists = getArtists(track.Artist)
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

func (mU *MusicUsecase) GetSongsByArtist(artist string) ([]models.VKTrack, error) {
	return mU.musicRepository.GetSongsByArtist(artist)
}

func (mU *MusicUsecase) GetSongById(id int) (models.VKTrack, error) {
	return mU.musicRepository.GetSongById(id)
}

func (mU *MusicUsecase) GetTracksCount() (int, error) {
	return mU.musicRepository.GetTracksCount()
}

func (mU *MusicUsecase) GetGenres() ([]string, error) {
	return mU.musicRepository.GetGenres()
}

func (mU *MusicUsecase) GetMusicByGenre(genre string) ([]models.VKTrack, error){
	return mU.musicRepository.GetMusicByGenre(genre)
}

func (mU *MusicUsecase) GetAllMusic() ([]models.VKTrack, error) {
	return mU.musicRepository.GetAllMusic()
}

func (mU *MusicUsecase) GetGenreFromHumanGenre(humanGenre string) (string, error) {
	return mU.musicRepository.GetGenreFromHumanGenre(humanGenre)
}

func (mU *MusicUsecase) GetArtistFromHumanArtist(humanArtist string) (string, error) {
	return mU.musicRepository.GetArtistFromHumanArtist(humanArtist)
}