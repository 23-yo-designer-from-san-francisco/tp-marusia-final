package repository

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"guessTheSongMarusia/models"
	log "guessTheSongMarusia/pkg/logger"
	"time"
)

const DAY = time.Hour * 24

type SessionRepository struct {
	redis *redis.Client
}

func NewSessionRepository(redis *redis.Client) *SessionRepository {
	return &SessionRepository{
		redis: redis,
	}
}

func (sR *SessionRepository) SaveSession(sessionID string, userSession *models.Session) error {
	sessionTillTime := time.Now().Add(DAY).Unix()
	expiration := time.Unix(sessionTillTime, 0)
	sessionJSON, err := json.Marshal(userSession)
	if err != nil {
		log.Error(err)
		return err
	}
	err = sR.redis.Set(sessionID, sessionJSON, time.Until(expiration)).Err()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (sR *SessionRepository) SavePlaylist(title string, tracks []models.VKTrack) error {
	tracksJSONE, err := json.Marshal(tracks)
	if err != nil {
		logrus.Error(err)
		return err
	}
	sessionTillTime := time.Now().Add(14 * DAY).Unix()
	expiration := time.Unix(sessionTillTime, 0)
	err = sR.redis.Set(title, tracksJSONE, time.Until(expiration)).Err()
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (sR *SessionRepository) GetPlaylist(title string) ([]models.VKTrack, error) {
	tracksJSON, err := sR.redis.Get(title).Result()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var tracks []models.VKTrack
	err = json.Unmarshal([]byte(tracksJSON), &tracks)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return tracks, nil
}

func (sR *SessionRepository) GetSession(sessionID string) (*models.Session, error) {
	sessionJSON, err := sR.redis.Get(sessionID).Result()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var userSession models.Session
	err = json.Unmarshal([]byte(sessionJSON), &userSession)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(userSession)
	return &userSession, nil
}

func (sR *SessionRepository) DeleteSession(sessionID string) error {
	err := sR.redis.Del(sessionID).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
