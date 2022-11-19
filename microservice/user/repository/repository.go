package repository

import (
	"guessTheSongMarusia/models"
	"time"
	log "guessTheSongMarusia/pkg/logger"
	"github.com/go-redis/redis"
	"encoding/json"
)

type SessionRepository struct {
	redis *redis.Client
}

func NewSessionRepository (redis *redis.Client) (*SessionRepository) {
	return &SessionRepository{
		redis: redis,
	}
}

func (sR *SessionRepository) SaveSession(sessionID string, userSession *models.Session) error {
	sessionTillTime:= time.Now().Add(time.Hour * 24).Unix()
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

func (sR *SessionRepository) DeleteSession(sessionID string) (error) {
	err := sR.redis.Del(sessionID).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}