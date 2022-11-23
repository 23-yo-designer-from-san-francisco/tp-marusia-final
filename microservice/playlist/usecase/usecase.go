package usecase

import (
	"errors"
	"guessTheSongMarusia/microservice/music"
	"guessTheSongMarusia/microservice/playlist"
	"guessTheSongMarusia/models"
	"guessTheSongMarusia/utils"
	"strings"

	"github.com/sirupsen/logrus"
)

const PlaylistAlreadyExistsErr = "playlist already exists"
const titleKeyLength = 32

type PlaylistUsecase struct {
	adjectives []string
	nouns []string
	playlistRepository playlist.Repository
	musicRepository music.Repository
}

func NewPlaylistUsecase(adjectives, nouns []string, playlistR playlist.Repository, musicR music.Repository) *PlaylistUsecase {
	return &PlaylistUsecase{
		adjectives: adjectives,
		nouns: nouns,
		playlistRepository: playlistR,
		musicRepository: musicR,
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

func (pU *PlaylistUsecase) GenerateTitle() (string, string, error) {
	for {
		title := utils.GeneratePlaylistName(pU.adjectives, pU.nouns)
		playlist, err := pU.playlistRepository.GetPlaylist(title)
		if err != nil && !strings.Contains(err.Error(), "redis: nil") {
			return "", "", err
		}
		if len(playlist) == 0 {
			titleKey := utils.RandStringBytesMaskImprSrcUnsafe(titleKeyLength)
			err := pU.playlistRepository.SaveTitle(titleKey, title)
			if err != nil {
				return "","", err
			}
			return titleKey, title, nil
		}
	}
}

func (pU *PlaylistUsecase) CreatePlaylist(playlist models.MiniAppPlaylist) (error) {
	title, err := pU.playlistRepository.GetTitle(playlist.TitleKey)
	if err != nil {
		return err
	}
	var VKTracks []models.VKTrack
	for _, musicId := range playlist.TracksIds {
		VKTrack, err := pU.musicRepository.GetMusicByID(musicId)
		if err != nil {
			return err
		}
		VKTrack.ArtistsWithHumanArtists, err = pU.musicRepository.GetArtistsInfoByMusicID(musicId)
		if err != nil {
			return err
		}
		VKTracks = append(VKTracks, *VKTrack)
	}
	err = pU.playlistRepository.SavePlaylist(title, VKTracks)
	if err != nil {
		return err
	}
	return nil
}