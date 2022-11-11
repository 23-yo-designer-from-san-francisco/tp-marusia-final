package repository

import (
	"github.com/jmoiron/sqlx"
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
	getTracksCount = `SELECT max(id) FROM music;`

	getGenres = `select title from genre;`
	getMusicByGenre = `select 
		m.title,
		m.artist,
		m.duration_two_url,
		m.duration_three_url,
		m.duration_five_url,
		m.duration_fifteen_url,
		m.human_title
		from music as m 
			join genre_music as gm on m.id = gm.music_id 
			join genre as g on g.id = gm.genre_id 
			where g.human_title = $1;`
	
	getAllSongs = `
		SELECT m.title,
				m.artist,
				m.duration_two_url,
				m.duration_three_url,
				m.duration_five_url,
				m.duration_fifteen_url,
				m.human_title
		FROM music AS m;
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

func (mR *MusicRepository) CreateTrack(track *models.VKTrack) error {
	tx, err := mR.db.Beginx()
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return err
	}

	var trackId int
	err = tx.QueryRowx(insertTrackQuery,
		&track.Title, &track.Artist,
		&track.Duration2, &track.Duration3, &track.Duration5, &track.Duration15,
		&track.HumanTitle).Scan(&trackId)

	if err != nil {
		log.Error(err)
		tx.Rollback()
		return err
	}
	log.Debug(trackId)
	for index, artist := range track.HumanArtists {
		_, err = tx.Exec(insertArtistQuery, &trackId, &track.Artists[index], &artist)
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (mR *MusicRepository) GetSongsByArtist(artist string) ([]models.VKTrack, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getSongsByHumanArtist, artist)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return VKTracks, nil
}

func (mR *MusicRepository) GetSongById(id int) (models.VKTrack, error) {
	var VKTrack models.VKTrack
	err := mR.db.Get(&VKTrack, getSongById, id)
	if err != nil {
		log.Error(err)
		return models.VKTrack{}, err
	}
	return VKTrack, nil
}

func (mR *MusicRepository) GetTracksCount() (int, error) {
	count := 0
	err := mR.db.Get(&count, getTracksCount)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	return count, nil
}

func (mR *MusicRepository) GetGenres() ([]string, error) {
	var genres = []string{}
	err := mR.db.Select(&genres, getGenres)
	if err != nil {
		log.Error(err)
		return []string{}, err
	}
	return genres, nil
}

func (mR *MusicRepository) GetMusicByGenre(genre string) ([]models.VKTrack, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getMusicByGenre, genre)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return VKTracks, nil
}

func (mR *MusicRepository) GetAllMusic() ([]models.VKTrack, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getAllSongs)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return VKTracks, nil
}