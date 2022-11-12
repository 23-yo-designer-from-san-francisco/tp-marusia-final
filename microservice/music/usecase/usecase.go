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

func validateMusicTrackTitle(title string) (string, bool) {
	title = html.UnescapeString(title)
	title = removeCharacters(title, ".,!?")
	//Удаляет скобки и содержимое
	reg := regexp.MustCompile(`\([^)]*\)`)
	title = reg.ReplaceAllString(title, "")
	title = strings.Replace(title, "$", "s", -1)
	title = strings.Replace(title, "é", "e", -1)
	//Если же что-то с чём-то, то нам такое не нужно
	return strings.ToLower(title), true
}

func getArtists(artist string) []string {
	var returnResult []string
	artists := strings.Split(artist, "feat.")
	for _, artistSplit := range artists {
		resultArtists := strings.Split(artistSplit, ",")
		for _, resultArtist := range resultArtists {
			if resultArtist[0] == ' ' {
				resultArtist = resultArtist[1:]
			}
			if resultArtist[len(resultArtist)-1] == ' ' {
				resultArtist = resultArtist[:len(resultArtist)-1]
			}
			returnResult = append(returnResult, resultArtist)
		}
	}
	return returnResult
}

func validateMusicArtists(artists []string) ([]string, bool) {
	var returnResult []string
	for _, artist := range artists {
		artist = html.UnescapeString(artist)
		artist = removeCharacters(artist, "!?")
		artist = strings.Replace(artist, "$", "s", -1)
		artist = strings.Replace(artist, "é", "e", -1)
		reg := regexp.MustCompile(`\([^)]*\)`)
		artist = reg.ReplaceAllString(artist, "")
		returnResult = append(returnResult, strings.ToLower(artist))
	}
	return returnResult, true
}

func validateMusicTrack(track *models.VKTrack) (*models.VKTrack, bool) {
	humanTitle, ok := validateMusicTrackTitle(track.Title)
	 if !ok {
	 	return nil, false
	 }
	track.HumanTitle = humanTitle
	track.Artists = getArtists(track.Artist)

	humanArtists, ok := validateMusicArtists(track.Artists)
	if !ok {
	 	return nil, false
	}
	track.HumanArtists = humanArtists
	return track, true
}

func (mU *MusicUsecase) CreateAllMusic(Tracks []models.VKTrack) error {
	for _, track := range Tracks {
		track, ok := validateMusicTrack(&track)
		if !ok {
			continue
		}
		log.Debug(track.Artists)
		log.Debug(track.HumanArtists)
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