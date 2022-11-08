package music

import "guessTheSongMarusia/models"

type Repository interface {
	CreateTrack(*models.VKTrack) (error)
}