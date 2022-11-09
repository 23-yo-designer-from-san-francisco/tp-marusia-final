package repository

import (
	//log "autfinal/pkg/logger"
	"guessTheSongMarusia/models"
	log "guessTheSongMarusia/pkg/logger"
	"github.com/jmoiron/sqlx"
)

const (
	insertTrackQuery = `insert into music (title, artist, duration_two_url, duration_three_url, duration_five_url, duration_fifteen_url, human_title)
		values ($1, $2, $3, $4, $5, $6, $7) returning id;`
	insertArtistQuery = `insert into artist (music_id, artist, human_artist) values ($1, $2, $3);`
	getSongsByHumanArtist =  `select m.title, m.artist, m.duration_two_url, m.duration_three_url, m.duration_five_url, m.duration_fifteen_url, m.human_title 
		from music as m join artist on artist.music_id = m.id where artist.human_artist = $1;
	`
)

type MusicRepository struct {
	db *sqlx.DB
}

func NewMusicRepository(db *sqlx.DB) *MusicRepository {
	return &MusicRepository{
		db: db,
	}
}

func(mR *MusicRepository) CreateTrack(track *models.VKTrack) (error) {
	tx, err := mR.db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}
	
	var trackId int 
	err = tx.QueryRowx(insertTrackQuery, 
		&track.Title, &track.Artist,
		&track.Duration2, &track.Duration3, &track.Duration5, &track.Duration15, 
		&track.HumanTitle).Scan(&trackId)

	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return err
	}
	log.Debug(trackId)
	for index, artist := range track.HumanArtists {
		_, err = tx.Exec(insertArtistQuery, &trackId, &track.Artists[index], &artist)
		if err != nil {
			log.Error(err.Error())
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