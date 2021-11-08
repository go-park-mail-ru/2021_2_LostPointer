package models

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
	"log"
)

type Track struct {
	ID          int64  `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Explicit    bool   `json:"explicit,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Number      int64  `json:"number,omitempty"`
	File        string `json:"file,omitempty"`
	ListenCount int64  `json:"listen_count,omitempty"`
	Duration    int64  `json:"duration,omitempty"`
	Lossless    bool   `json:"lossless,omitempty"`
	Album       Album  `json:"album"`
	Artist      Artist `json:"artist"`
}

func (t *Track) BindProtoTrack(track *proto.Track) {
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
		Album: Album{
			ID:             track.Album.ID,
			Title:          track.Album.Title,
			Artwork:        track.Album.Artwork,
			Year:           track.Album.Year,
			Artist:         track.Album.Artist,
			TracksCount:    track.Album.TracksAmount,
			TracksDuration: track.Album.TracksDuration,
		},
		Artist: Artist{
			ID:     track.Artist.ID,
			Name:   track.Artist.Name,
			Avatar: track.Artist.Avatar,
			Video:  track.Artist.Video,
		},
	}

	log.Println(bindedTrack)

	*t = *bindedTrack
}
