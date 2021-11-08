package models

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
)

type Track struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Explicit    bool   `json:"explicit"`
	Genre       string `json:"genre,omitempty"`
	Number      int64  `json:"number,omitempty"`
	File        string `json:"file"`
	ListenCount int64  `json:"listenCount"`
	Duration    int64  `json:"duration"`
	Lossless    bool   `json:"lossless"`
	Album       Album  `json:"album"`
	Artist      Artist `json:"artist,omitempty"`
}

func (t Track) BindProtoTrack(track *proto.Track) *Track {
	bindedTrack := &Track{
		ID:          track.ID,
		Title:       track.Title,
		Explicit:    track.Explicit,
		Genre:       track.Genre,
		Number:      track.Number,
		File:        track.File,
		ListenCount: track.ListenCount,
		Duration:    track.Duration,
		Lossless:    track.Lossless,
		Album:       Album{
			ID: track.Album.ID,
			Title: track.Album.Title,
			Artwork: track.Album.Artwork,
		},
		Artist:      Artist{
			ID: track.Artist.ID,
			Name: track.Artist.Name,
		},
	}

	return bindedTrack
}
