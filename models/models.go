package models

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"guessTheSongMarusia/utils"
)

type Duration int64

const (
	Three  Duration = 3
	Five   Duration = 5
	Eight  Duration = 8
	Thirty Duration = 30
)

const (
	GuessedAttempt1 = 12
	GuessedAttempt2 = 8
	GuessedAttempt3 = 4
)

const (
	StatusNewGame = iota
	StatusChoosingGenre
	StatusListingGenres
	StatusPlaying
	StatusNewCompetition
	StatusCompetitionRules
	StatusCompetition
	StatusChooseArtist
	StatusGeneratedPlaylist
	StatusPlaylistFinished
	StatusKeyPhrasePlaylist
)

const (
	GenreMode = iota
	ArtistMode
)

const LevensteinSimilarityLimit = 0.85

type Track struct {
	Id     int64
	Title  string
	Artist string
	Audio  map[Duration]string
}

type Session struct {
	CurrentLevel    Duration
	CurrentPoints   int
	CurrentTrack    VKTrack
	MusicStarted    bool
	NextLevelLoses  bool
	TitleMatch      bool
	ArtistMatch     bool
	PlayedTracks    map[int]bool
	CurrentGenre    string
	GameMode        int
	TrackCounter    int
	GameState       *State
	CurrentPlaylist []VKTrack
	CompetitionMode bool
	KeyPhrase       string
}

func NewSession() *Session {
	return &Session{
		PlayedTracks:  make(map[int]bool),
		GameState:     NewGameState,
		GameMode:      GenreMode,
		CurrentPoints: 0,
	}
}

type TracksPerGenres struct {
	Rock    []VKTrack `json:"rock"`
	NotRock []VKTrack `json:"not_rock"`
}

type VKTrack struct {
	ID                      int                 `json:"id,omitempty" db:"id"`
	Title                   string              `json:"title,omitempty" db:"title"`
	Artist                  string              `json:"artist,omitempty" db:"artist"`
	Duration3               string              `json:"duration_3,omitempty" db:"duration_three_url"`
	Duration5               string              `json:"duration_5,omitempty" db:"duration_five_url"`
	Duration8               string              `json:"duration_8,omitempty" db:"duration_eight_url"`
	Duration30              string              `json:"duration_30,omitempty" db:"duration_thirty_url"`
	ArtistsWithHumanArtists map[string][]string `json:"human_artists"`
	HumanTitles             pq.StringArray      `json:"human_titles" db:"human_titles"`
}

type MiniAppTrack struct {
	ID     int    `json:"id,omitempty" db:"id"`
	Title  string `json:"title,omitempty" db:"title"`
	Artist string `json:"artist,omitempty" db:"artist"`
}

type TitleInfo struct {
	TitleKey string `json:"title_key,omitempty"`
	Title    string `json:"title,omitempty"`
}

type MiniAppPlaylist struct {
	TitleKey  string `json:"title_key,omitempty"`
	TracksIds []int  `json:"music_ids,omitempty"`
}

func (track *VKTrack) checkTitleInAnswer(answer string) bool {
	if !(utils.ContainsAny(answer, track.Title) && utils.ContainsAny(answer, track.HumanTitles...)) {
		log.Debug("Checking by Levenshtein")
		lev := metrics.NewLevenshtein()
		for _, title := range track.HumanTitles {
			similarity := strutil.Similarity(answer, title, lev)
			log.WithFields(log.Fields{
				"limit":      LevensteinSimilarityLimit,
				"similarity": similarity,
				"answer":     answer,
				"title":      title,
			}).Debug("Debug similarity")
			if similarity >= LevensteinSimilarityLimit {
				return true
			}
		}
		return false
	}
	return true
}

func (track *VKTrack) checkArtistsInAnswer(answer string) bool {
	for artistName, humanNames := range track.ArtistsWithHumanArtists {
		if utils.ContainsAny(answer, artistName) || utils.ContainsAny(answer, humanNames...) {
			return true
		}
	}
	log.Debug("Checking by Levenshtein")
	lev := metrics.NewLevenshtein()
	for _, humanNames := range track.ArtistsWithHumanArtists {
		for _, name := range humanNames {
			similarity := strutil.Similarity(answer, name, lev)
			log.WithFields(log.Fields{
				"limit":      LevensteinSimilarityLimit,
				"similarity": similarity,
				"answer":     answer,
				"name":       name,
			}).Debug("Debug similarity")
			if similarity >= LevensteinSimilarityLimit {
				return true
			}
		}
	}
	return false
}

func (track *VKTrack) CheckUserAnswer(answer string, userSession *Session) (bool, bool) {
	matchTitle := track.checkTitleInAnswer(answer)
	matchArtist := track.checkArtistsInAnswer(answer)
	if matchTitle {
		userSession.TitleMatch = true
	}
	if matchArtist {
		userSession.ArtistMatch = true
	}
	return matchTitle, matchArtist
}

type VKTracks struct {
	VKTracks []VKTrack
}
