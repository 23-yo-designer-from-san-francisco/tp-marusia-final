package models

type Duration int64
type GameStatus int32

const (
	Two  Duration = 2
	Five Duration = 5
	Ten  Duration = 10
)

const (
	New           GameStatus = 0
	ChoosingGenre GameStatus = 1
	ListingGenres GameStatus = 2
	Playing       GameStatus = 3
)

type Track struct {
	Id     int64
	Title  string
	Artist string
	Audio  map[Duration]string
}

type Session struct {
	CurrentLevel   Duration
	CurrentPoints  int64
	CurrentTrack   VKTrack
	GameStatus     GameStatus
	MusicStarted   bool
	NextLevelLoses bool
	TitleMatch     bool
	ArtistMatch    bool
	PlayedTracks   map[int]bool
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
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Duration2 string `json:"duration_2"`
	Duration3 string `json:"duration_3"`
	Duration5 string `json:"duration_5"`
}
