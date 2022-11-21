package music

import "guessTheSongMarusia/models"

type Usecase interface {
	CreateAllMusic(tracks []models.VKTrack) error
	GetSongsByArtist(artist string) ([]models.VKTrack, string, error)
	GetGenres() ([]string, error)
	GetMusicByGenre(genre string) ([]models.VKTrack, error)
	GetAllMusic() ([]models.VKTrack, error)
	GetGenreFromHumanGenre(humanGenre string) (string, error)
	GetRandomMusic(limit int) ([]models.VKTrack, error)
	GetRandomMusicByGenre(limit int, humanGenre string) ([]models.VKTrack, error)
	GetRandomMusicByArtist(limit int, humanArtist string) ([]models.VKTrack, string, error)
}
