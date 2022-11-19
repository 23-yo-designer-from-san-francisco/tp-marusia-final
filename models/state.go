package models

type phraseFunc func() (string, string)

type State struct {
	GameStatus int
	StandartPhraseText string
	StandartPhraseTTS string
}

func NewState(gameStatus int, phraseF, errorPhraseF phraseFunc) *State{
	state := State{}
	state.GameStatus = gameStatus
	state.StandartPhraseText, state.StandartPhraseTTS = phraseF()
	return &state
}

func (state *State) SayStandartPhrase() (string, string) {
	return state.StandartPhraseText, state.StandartPhraseTTS
}

func (state *State) SayErrorPhrase() (string, string) {
	return state.StandartPhraseText, state.StandartPhraseTTS
}

var NewGameState = NewState(StatusNewGame, StartGamePhrase, StandartErrorPhrase) 
var ChooseGenreState = NewState(StatusChoosingGenre, ChooseGenrePhrase, StandartErrorPhrase)
var ListingGenreState = NewState(StatusListingGenres, AvailableGenresPhrase, StandartErrorPhrase)
var PlayingState = NewState(StatusPlaying, PlayingGamePhrase, StandartErrorPhrase)
var NewCompetitionState = NewState(StatusNewCompetition, CompetitionPhrase, StandartErrorPhrase)
var CompetitonRulesState = NewState(StatusCompetitionRules, CompetitionRulesPhrase, StandartErrorPhrase)
var ChooseArtistState = NewState(StatusChooseArtist, ChooseArtistPhrase, StandartErrorPhrase)
