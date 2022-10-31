package models

type Duration int64

const (
	TwoSecondsLevel = iota
	FiveSecondsLevel
	TenSecondsLevel
)

const (
	Two  Duration = 2
	Five Duration = 5
	Ten  Duration = 10
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
	GameStarted    bool
	MusicStarted   bool
	NextLevelLoses bool
	PlayedTracks   map[int]bool
}

func NewSession() *Session {
	return &Session{
		PlayedTracks: make(map[int]bool),
	}
}

type VKTrack struct {
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Duration2 string `json:"duration_2"`
	Duration3 string `json:"duration_3"`
	Duration5 string `json:"duration_5"`
}
