package repository

import (
	"guessTheSongMarusia/models"
	log "guessTheSongMarusia/pkg/logger"
	"strings"

	"github.com/jmoiron/sqlx"
)

// https://sqlformat.org/
const (
	insertMusicQuery = `
		INSERT INTO music (title, artist, duration_two_url, duration_three_url, duration_five_url, duration_fifteen_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;`

	selectArtistIDQuery = `select id from artist where artist = $1;`

	insertArtistQuery     = `insert into artist (music_id, artist, human_artist) values ($1, $2, $3);`
	insertArtistV2Query   = `insert into artist (music_id, artist) values ($1,$2) returning id;`

	insertArtistMusic = `insert into artist_music (music_id, artist_id) values ($1, $2);`
	insertHumanTitle = `insert into human_title (music_id, human_title) values ($1, $2);`
	insertHumanArtist = `insert into human_artist (artist_id, human_artist) values ($1, $2);`

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
	
	getGenreFromHumanGenre = `select title from genre where human_title = $1;`

	getArtistFromHumanArtist = `select distinct artist from artist where human_artist = $1`;
	
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

func (mR *MusicRepository) GetGenreFromHumanGenre(humanGenre string) (string, error) {
	var genre string
	err := mR.db.Get(&genre, getGenreFromHumanGenre, humanGenre)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return genre, nil
}

func (mR *MusicRepository) GetArtistFromHumanArtist(humanArtist string) (string, error) {
	var artist string
	err := mR.db.Get(&artist, getArtistFromHumanArtist, humanArtist)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return artist, nil
}

func (mR *MusicRepository) CreateMusic(track *models.VKTrack) (error) {
	tx, err := mR.db.Beginx()
	if err != nil {
		return err
	}
	var musicID int
	err = tx.QueryRowx(insertMusicQuery, 
		&track.Title, 
		&track.Artist, 
		&track.Duration2, 
		&track.Duration3, 
		&track.Duration5, 
		&track.Duration15).Scan(&musicID)

	//Значит, такой трек уже есть
	if err != nil {
		log.Error("Значит трек есть:", err)
		tx.Rollback()
		return err
	}

	//Записываем артистов и человеческих артистов и берём id
	var ArtistIds []int
	for artist, humanArtistNames := range track.ArtistsWithHumanArtists {
		var artistID int
		err := tx.QueryRowx(selectArtistIDQuery, &artist).Scan(&artistID)
		//Если артиста нет
		if err != nil {
			if !strings.Contains(err.Error(), "no rows in result set") {
				log.Error("Смотрим ошибку после селекта ", err)
				tx.Rollback()
				return err
			}

			err := tx.QueryRowx(insertArtistV2Query, &musicID, &artist).Scan(&artistID)
			if err != nil {
				log.Error("Тут на ошибку с базой данных реагируем артист инсерт ", err)
				tx.Rollback()
				return err
			}

			for _, humanArtist := range humanArtistNames {
				_, err = tx.Exec(insertHumanArtist, &artistID, humanArtist)
				if err != nil {
					log.Error("human_artist", err)
					tx.Rollback()
					return err
				}
			}
		}
		ArtistIds = append(ArtistIds, artistID)
	}

	//Заполняем много ко многим артист с музыкой
	for _, artistID := range ArtistIds {
		_, err := tx.Exec(insertArtistMusic, &musicID, &artistID)
		if err != nil {
			log.Error("Тут на ошибку с базой данных реагируем", err)
			tx.Rollback()
			return err
		}
	}

	//Заполняем человеческие названия по умолчанию(Базовая валидация)
	for _, humanTitle := range track.HumanTitles {
		_, err = tx.Exec(insertHumanTitle, &musicID, &humanTitle)
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return err
		}
	} 
	
	tx.Commit()
	return nil;
}