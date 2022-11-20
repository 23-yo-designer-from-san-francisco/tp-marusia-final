package usecase

import (
	"guessTheSongMarusia/microservice/user"
	"guessTheSongMarusia/models"
)

type sessionUsecase struct {
	SessionRepo user.SessionRepository
}

func NewUserUsecase(sessionRepo user.SessionRepository) *sessionUsecase {
	return &sessionUsecase{
		SessionRepo: sessionRepo,
	}
}

func (sU *sessionUsecase) SaveSession(session string, userSession *models.Session) error {
	return sU.SessionRepo.SaveSession(session, userSession)
}

func (sU *sessionUsecase) GetSession(session string) (*models.Session, error) {
	return sU.SessionRepo.GetSession(session)
}

func (sU *sessionUsecase) DeleteSession(session string) error {
	return sU.SessionRepo.DeleteSession(session)
}

func (sU *sessionUsecase) GetPlaylist(title string) ([]models.VKTrack, error) {
	return sU.SessionRepo.GetPlaylist(title)
}

func (sU *sessionUsecase) SavePlaylist(title string, tracks []models.VKTrack) error {
	return sU.SessionRepo.SavePlaylist(title, tracks)
}
