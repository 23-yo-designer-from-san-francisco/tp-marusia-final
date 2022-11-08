package music

import "guessTheSongMarusia/models"

type Usecase interface {
	CreateAllMusic(tracks []models.VKTrack) error
}