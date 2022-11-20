package music

import "guessTheSongMarusia/models"

type Repository interface {
	GetSongsByArtist(artist string) ([]models.VKTrack, string, error)
	GetGenres() ([]string, error)
	GetMusicByGenre(genre string) ([]models.VKTrack, error)
	GetAllMusic() ([]models.VKTrack, error)
	GetGenreFromHumanGenre(humanGenre string) (string, error)
	CreateTrack(track *models.VKTrack) (error)
	GetRandomMusic(limit int) ([]models.VKTrack, error)
	GetRandomMusicByGenre(limit int, humanGenre string) ([]models.VKTrack, error)
}
