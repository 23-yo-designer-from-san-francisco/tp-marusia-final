package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	//log "autfinal/pkg/logger"
	"guessTheSongMarusia/models"
	log "guessTheSongMarusia/pkg/logger"
)

// https://sqlformat.org/
const (
	insertTrackQuery = `
		INSERT INTO music (title, artist, duration_two_url, duration_three_url, duration_five_url, duration_fifteen_url, human_title)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
`
	insertArtistQuery     = `insert into artist (music_id, artist, human_artist) values ($1, $2, $3);`
	getSongsByHumanArtist = `
		SELECT m.title,
			   m.artist,
			   m.duration_two_url,
			   m.duration_three_url,
			   m.duration_five_url,
			   m.duration_fifteen_url,
			   m.human_title
		FROM music AS m
		JOIN artist ON artist.music_id = m.id
		WHERE artist.human_artist = $1;
	`
	getSongById = `
		SELECT m.title,
			   m.artist,
			   m.duration_two_url,
			   m.duration_three_url,
			   m.duration_five_url,
			   m.duration_fifteen_url,
			   m.human_title
		FROM music AS m
		JOIN artist ON artist.music_id = m.id
		WHERE m.id = $1;
	`
	getTracksCount = `SELECT max(id) FROM music`
)

type MusicRepository struct {
	db *sqlx.DB
}

func NewMusicRepository(db *sqlx.DB) *MusicRepository {
	return &MusicRepository{
		db: db,
	}
}

func (mR *MusicRepository) CreateTrack(track *models.VKTrack) error {
	tx, err := mR.db.Beginx()
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	var trackId int
	err = tx.QueryRowx(insertTrackQuery,
		&track.Title, &track.Artist,
		&track.Duration2, &track.Duration3, &track.Duration5, &track.Duration15,
		&track.HumanTitle).Scan(&trackId)

	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		log.Error(err.Error())
		return err
	}
	log.Debug(trackId)
	for index, artist := range track.HumanArtists {
		_, err = tx.Exec(insertArtistQuery, &trackId, &track.Artists[index], &artist)
		if err != nil {
			logrus.Error(err.Error())
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (mR *MusicRepository) GetSongsByArtists(artist string) ([]models.VKTrack, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getSongsByHumanArtist, artist)
	if err != nil {
		return nil, err
	}
	return VKTracks, nil
}

func (mR *MusicRepository) GetSongById(id int) (models.VKTrack, error) {
	var VKTracks []models.VKTrack
	err := mR.db.Select(&VKTracks, getSongById, id)
	if err != nil {
		return models.VKTrack{}, err
	}
	return VKTracks[0], nil
}

func (mR *MusicRepository) GetTracksCount() (int, error) {
	count := []int{0}
	err := mR.db.Select(&count, getTracksCount)

	if err != nil {
		return -1, err
	}
	return count[0], nil
}
