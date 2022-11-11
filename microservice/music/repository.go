package music

import "guessTheSongMarusia/models"

type Repository interface {
	CreateTrack(*models.VKTrack) error
	GetSongsByArtist(artist string) ([]models.VKTrack, error)
	GetSongById(id int) (models.VKTrack, error)
	GetTracksCount() (int, error)
	GetGenres() ([]string, error)
	GetMusicByGenre(genre string) ([]models.VKTrack, error)
	GetAllMusic() ([]models.VKTrack, error)
}
