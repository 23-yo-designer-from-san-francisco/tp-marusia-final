package playlist

import "guessTheSongMarusia/models"

type Repository interface {
	SavePlaylist(title string, tracks []models.VKTrack) error
	GetPlaylist(title string) ([]models.VKTrack, error)
	SaveTitle(titleKey, title string) (error)
	GetTitle(titleKey string) (string, error)
}
