package models

type phraseFunc func() (string, string)

type State struct {
	GameStatus int
	StandartPhraseText string
	StandartPhraseTTS string
}

func NewState(gameStatus int, phraseF phraseFunc) *State{
	state := State{}
	state.GameStatus = gameStatus
	state.StandartPhraseText, state.StandartPhraseTTS = phraseF()
	return &state
}

func (state *State) SayStandartPhrase() (string, string) {
	return state.StandartPhraseText, state.StandartPhraseTTS
}

var NewGameState = NewState(StatusNewGame, StartGamePhrase) 
var ChooseGenreState = NewState(StatusChoosingGenre, ChooseGenrePhrase)
var ListingGenreState = NewState(StatusListingGenres, AvailableGenresPhrase)
var PlayingState = NewState(StatusPlaying, PlayingGamePhrase)
var NewCompetitionState = NewState(StatusNewCompetition, CompetitionPhrase)
var CompetitonRulesState = NewState(StatusCompetitionRules, CompetitionRulesPhrase)
var ChooseArtistState = NewState(StatusChooseArtist, ChooseArtistPhrase)
