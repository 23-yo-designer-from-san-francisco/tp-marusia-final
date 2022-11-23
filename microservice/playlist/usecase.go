package playlist

import "guessTheSongMarusia/models"

type Usecase interface {
	SavePlaylist(title string, tracks []models.VKTrack) error
	GetPlaylist(title string) ([]models.VKTrack, error)
}
