package usecase

import (
	"errors"
	"github.com/sirupsen/logrus"
	"guessTheSongMarusia/microservice/playlist"
	"guessTheSongMarusia/models"
)

const PlaylistAlreadyExistsErr = "playlist already exists"

type PlaylistUsecase struct {
	playlistRepository playlist.Repository
}

func NewPlaylistUsecase(playlistR playlist.Repository) *PlaylistUsecase {
	return &PlaylistUsecase{
		playlistRepository: playlistR,
	}
}

func (pU *PlaylistUsecase) GetPlaylist(title string) ([]models.VKTrack, error) {
	return pU.playlistRepository.GetPlaylist(title)
}

func (pU *PlaylistUsecase) SavePlaylist(title string, tracks []models.VKTrack) error {
	tracks, err := pU.playlistRepository.GetPlaylist(title)
	// Если плейлист ЕСТЬ, ошибки НЕТ
	logrus.Warn("SavePlaylist err: ", err)
	if err == nil {
		return errors.New(PlaylistAlreadyExistsErr)
	}
	return pU.playlistRepository.SavePlaylist(title, tracks)
}
