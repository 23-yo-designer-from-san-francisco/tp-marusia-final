package user

import "guessTheSongMarusia/models"

type SessionUsecase interface {
	SaveSession(session string, userSession *models.Session) error
	GetSession(sessionID string) (*models.Session, error)
	DeleteSession(sessionID string) error
}
