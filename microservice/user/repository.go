package user

import "guessTheSongMarusia/models"

type SessionRepository interface {
	SaveSession(session string, userSession *models.Session) error
	GetSession(sessionID string) (*models.Session, error)
	DeleteSession(sessionID string) error
	SavePlaylist(title string, tracks []models.VKTrack) error
	GetPlaylist(title string) ([]models.VKTrack, error)
}
