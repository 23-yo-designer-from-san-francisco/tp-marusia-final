package music

import "guessTheSongMarusia/models"

type Usecase interface {
	CreateAllMusic(tracks []models.VKTrack) error
	GetSongsByArtists(artist string) ([]models.VKTrack, error)
	GetSongById(id int) (models.VKTrack, error)
	GetTracksCount() (int, error)
}
