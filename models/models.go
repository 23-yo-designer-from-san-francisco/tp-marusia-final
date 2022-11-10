package models

type Duration int64

const (
	Two  Duration = 2
	Five Duration = 5
	Ten  Duration = 10
)

const (
	New = iota
	ChoosingGenre
	ListingGenres
	Playing
	CompetitionIntro
	CompetitionRules
	Competition
)

type Track struct {
	Id     int64
	Title  string
	Artist string
	Audio  map[Duration]string
}

type Session struct {
	CurrentLevel      Duration
	CurrentPoints     int64
	CurrentTrack      VKTrack
	GameStatus        int
	MusicStarted      bool
	NextLevelLoses    bool
	TitleMatch        bool
	ArtistMatch       bool
	PlayedTracks      map[int]bool
	CurrentGenre      string
	GenreTrackCounter int
	Fails             int
}

func NewSession() *Session {
	return &Session{
		PlayedTracks: make(map[int]bool),
		GameStatus:   New,
	}
}

type TracksPerGenres struct {
	Rock    []VKTrack `json:"rock"`
	NotRock []VKTrack `json:"not_rock"`
}

type VKTrack struct {
	Title        string   `json:"title,omitempty" db:"title"`
	Artist       string   `json:"artist,omitempty" db:"artist"`
	Duration2    string   `json:"duration_2,omitempty" db:"duration_two_url"`
	Duration3    string   `json:"duration_3,omitempty" db:"duration_three_url"`
	Duration5    string   `json:"duration_5,omitempty" db:"duration_five_url"`
	Duration15   string   `json:"duration_15,omitempty" db:"duration_fifteen_url"`
	Artists      []string `json:"-"`
	HumanTitle   string   `json:"human_title" db:"human_title"`
	HumanArtists []string `json:"-"`
}

type VKTracks struct {
	VKTracks []VKTrack
}
