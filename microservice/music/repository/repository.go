package repository

import (
	"guessTheSongMarusia/models"
	log "guessTheSongMarusia/pkg/logger"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// https://sqlformat.org/
const (
	insertMusicQueryV2 = `
		INSERT INTO music (title, artist, duration_three_url, duration_five_url, duration_eight_url, duration_thirty_url, human_titles)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	selectArtistIDQuery = `select id from artist where artist_name = $1;`

	insertArtistQueryV2 = `insert into artist (artist_name, human_artists) values ($1, $2) returning id;`

	insertArtistMusic = `insert into artist_music (music_id, artist_id) values ($1, $2);`

	selectArtistsInfoByMusicId = `
		select a.artist_name, a.human_artists from artist_music as am
		join artist as a on a.id = am.artist_id
		where am.music_id = $1;`

	getSongsByHumanArtist = `select 
		m.id,
		m.title,
		m.artist,
		m.duration_three_url,
		m.duration_five_url,
		m.duration_eight_url,
		m.duration_thirty_url,
		m.human_titles from music as m
		join artist_music as am on m.id = am.music_id 
		join artist as a on a.id = am.artist_id
		where $1 = ANY(a.human_artists);`

	getGenres = `select genre from genre;`

	getMusicByGenre = `select 
    	m.id,
		m.title,
		m.artist,
		m.duration_three_url,
		m.duration_five_url,
		m.duration_eight_url,
		m.duration_thirty_url,
		m.human_titles
		from music as m 
			join genre_music as gm on m.id = gm.music_id 
			join genre as g on g.id = gm.genre_id 
			where $1 = ANY(g.human_genres);`

	getArtistFromHumanArtist = `select distinct artist_name from artist where $1 = ANY(human_artists)`
	getGenreFromHumanGenre   = `select genre from genre where $1 = ANY(human_genres);`

	getAllSongs = `
		SELECT 
				m.id,
				m.title,
				m.artist,
				m.duration_three_url,
				m.duration_five_url,
				m.duration_eight_url,
				m.duration_thirty_url,
				m.human_titles
		FROM music AS m;
	`
	getRandomMusic = `
		SELECT
				m.id,
				m.title,
				m.artist,
				m.duration_three_url,
				m.duration_five_url,
				m.duration_eight_url,
				m.duration_thirty_url,
				m.human_titles
		FROM music AS m
		ORDER BY random()
		LIMIT $1;`
	
	getRandomMusicByGenre = `
		select 
    	m.id,
		m.title,
		m.artist,
		m.duration_three_url,
		m.duration_five_url,
		m.duration_eight_url,
		m.duration_thirty_url,
		m.human_titles
		from music as m 
			join genre_music as gm on m.id = gm.music_id 
			join genre as g on g.id = gm.genre_id 
			where $1 = ANY(g.human_genres)
			order by random()
			limit $2;`
	
	getRandomMusicByArtist = `
		select 
    	m.id,
		m.title,
		m.artist,
		m.duration_three_url,
		m.duration_five_url,
		m.duration_eight_url,
		m.duration_thirty_url,
		m.human_titles
		from music as m 
			join artist_music as am on m.id = am.music_id 
			join artist as a on a.id = am.artist_id 
			where $1 = ANY(a.human_artists)
			order by random()
			limit $2;`

	getMusicByID = `select * from music where id = $1;`
)

type MusicRepository struct {
	db *sqlx.DB
}

func NewMusicRepository(db *sqlx.DB) *MusicRepository {
	return &MusicRepository{
		db: db,
	}
}

func(mR *MusicRepository) fillTracksWithArtists(VKTracks []models.VKTrack) ([]models.VKTrack, error) {
	var err error
	for index, track := range VKTracks {
		VKTracks[index].ArtistsWithHumanArtists, err = mR.GetArtistsInfoByMusicID(track.ID)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return VKTracks, nil
}

func (mR *MusicRepository) GetMusicByID(musicId int) (*models.VKTrack, error) {
	var track models.VKTrack
	err := mR.db.Get(&track, getMusicByID, &musicId)
	if err != nil {
		return nil, err
	}
	return &track, nil
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

func (mR *MusicRepository) GetGenreFromHumanGenre(humanGenre string) (string, error) {
	var genre string
	err := mR.db.Get(&genre, getGenreFromHumanGenre, humanGenre)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return genre, nil
}

func (mR *MusicRepository) CreateTrack(track *models.VKTrack) error {
	tx, err := mR.db.Beginx()
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return err
	}

	var musicID int
	err = tx.QueryRowx(insertMusicQueryV2,
		&track.Title,
		&track.Artist,
		&track.Duration3,
		&track.Duration5,
		&track.Duration8,
		&track.Duration30,
		&track.HumanTitles).Scan(&musicID)

	if err != nil {
		//Если такой трек есть уже или же база упала можем выходить
		log.Error(err)
		tx.Rollback()
		return err
	}

	for artist, humanArtists := range track.ArtistsWithHumanArtists {
		var artistID int
		err := tx.Get(&artistID, selectArtistIDQuery, artist)
		if err != nil {
			//Проверка, не упала ли база
			if !strings.Contains(err.Error(), "no rows in result set") {
				log.Error("Смотрим ошибку1", err)
				tx.Rollback()
				return err
			}
			//Если такого артиста нет
			err = tx.QueryRow(insertArtistQueryV2,
				&artist,
				&humanArtists).Scan(&artistID)

			if err != nil {
				log.Error("Смотрим ошибку2", err)
				tx.Rollback()
				return err
			}
		}
		//Заполняем связь много ко многим
		_, err = tx.Exec(insertArtistMusic, &musicID, &artistID)
		if err != nil {
			log.Error("Смотрим ошибку3", err)
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (mR *MusicRepository) GetArtistsInfoByMusicID(musicID int) (map[string][]string, error) {
	type artistInfo struct {
		Artist       string         `db:"artist_name"`
		HumanArtists pq.StringArray `db:"human_artists"`
	}

	var artistStructs []artistInfo
	err := mR.db.Select(&artistStructs, selectArtistsInfoByMusicId, &musicID)
	if err != nil {
		log.Error("Смотрим ошибку4", err)
		return nil, err
	}
	artistsMap := make(map[string][]string)
	for _, artist := range artistStructs {
		log.Debug(artist)
		artistsMap[artist.Artist] = artist.HumanArtists
	}
	return artistsMap, nil
}

func (mR *MusicRepository) GetAllMusic() ([]models.VKTrack, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getAllSongs)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	VKTracks, err = mR.fillTracksWithArtists(VKTracks)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return VKTracks, nil
}

func (mR *MusicRepository) GetMusicByGenre(genre string) ([]models.VKTrack, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getMusicByGenre, genre)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	VKTracks, err = mR.fillTracksWithArtists(VKTracks)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return VKTracks, nil
}

func (mR *MusicRepository) GetSongsByArtist(human_artist string) ([]models.VKTrack, string, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getSongsByHumanArtist, human_artist)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}
	var artist string
	err = mR.db.Get(&artist, getArtistFromHumanArtist, human_artist)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}

	VKTracks, err = mR.fillTracksWithArtists(VKTracks)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}
	return VKTracks, artist, nil
}

func (mR *MusicRepository) GetRandomMusic(limit int) ([]models.VKTrack, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getRandomMusic, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	VKTracks, err = mR.fillTracksWithArtists(VKTracks)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return VKTracks, nil
}

func (mR *MusicRepository) GetRandomMusicByGenre(limit int, humanGenre string) ([]models.VKTrack, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getRandomMusicByGenre, humanGenre, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	VKTracks, err = mR.fillTracksWithArtists(VKTracks)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return VKTracks, nil
}

func (mR *MusicRepository) GetRandomMusicByArtist(limit int, humanArtist string) ([]models.VKTrack, string, error) {
	var VKTracks = []models.VKTrack{}
	err := mR.db.Select(&VKTracks, getRandomMusicByArtist, humanArtist, limit)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}

	var artist string
	err = mR.db.Get(&artist, getArtistFromHumanArtist, humanArtist)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}

	VKTracks, err = mR.fillTracksWithArtists(VKTracks)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}
	return VKTracks, artist, nil
}