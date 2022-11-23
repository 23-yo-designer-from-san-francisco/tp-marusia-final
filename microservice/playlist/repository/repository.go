package repository

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"guessTheSongMarusia/models"
	log "guessTheSongMarusia/pkg/logger"
	"time"
)

type PlaylistRepository struct {
	//db *sqlx.DB
	redis *redis.Client
}

func NewPlaylistRepository( /*db *sqlx.DB*/ redis *redis.Client) *PlaylistRepository {
	//return &PlaylistRepository{
	//	db: db,
	//}
	return &PlaylistRepository{
		redis: redis,
	}
}

func (sR *PlaylistRepository) SavePlaylist(title string, tracks []models.VKTrack) error {
	tracksJSONE, err := json.Marshal(tracks)
	if err != nil {
		logrus.Error(err)
		return err
	}
	sessionTillTime := time.Now().Add(14 * models.DAY).Unix()
	expiration := time.Unix(sessionTillTime, 0)
	err = sR.redis.Set(title, tracksJSONE, time.Until(expiration)).Err()
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (pR *PlaylistRepository) GetPlaylist(title string) ([]models.VKTrack, error) {
	tracksJSON, err := pR.redis.Get(title).Result()
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
