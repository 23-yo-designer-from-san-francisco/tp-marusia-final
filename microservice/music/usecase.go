package music

import "guessTheSongMarusia/models"

type Usecase interface {
	CreateAllMusic(tracks []models.VKTrack) error
	GetSongsByArtist(artist string) ([]models.VKTrack, error)
	GetGenres() ([]string, error)
	GetMusicByGenre(genre string) ([]models.VKTrack, error)
	GetAllMusic() ([]models.VKTrack, error)
	GetGenreFromHumanGenre(humanGenre string) (string, error)
}
