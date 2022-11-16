package music

import "guessTheSongMarusia/models"

type Repository interface {
	GetSongsByArtist(artist string) ([]models.VKTrack, error)
	GetGenres() ([]string, error)
	GetMusicByGenre(genre string) ([]models.VKTrack, error)
	GetAllMusic() ([]models.VKTrack, error)
	GetGenreFromHumanGenre(humanGenre string) (string, error)
	CreateTrack(track *models.VKTrack) (error)
}
