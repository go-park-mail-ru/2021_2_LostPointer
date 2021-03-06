package models

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
)

//easyjson:json
type (
	Tracks []Track

	Track struct {
		ID            int64  `json:"id,omitempty"`
		Title         string `json:"title,omitempty"`
		Explicit      bool   `json:"explicit,omitempty"`
		Genre         string `json:"genre,omitempty"`
		Number        int64  `json:"number,omitempty"`
		File          string `json:"file,omitempty"`
		ListenCount   int64  `json:"listen_count,omitempty"`
		Duration      int64  `json:"duration,omitempty"`
		Lossless      bool   `json:"lossless,omitempty"`
		Album         Album  `json:"album,omitempty"`
		Artist        Artist `json:"artist,omitempty"`
		IsInFavorites bool   `json:"is_in_favorites,omitempty"`
	}

	TrackAlbum struct {
		ID            int64  `json:"id,omitempty"`
		Title         string `json:"title,omitempty"`
		Explicit      bool   `json:"explicit"`
		Genre         string `json:"genre,omitempty"`
		Number        int64  `json:"number,omitempty"`
		File          string `json:"file,omitempty"`
		ListenCount   int64  `json:"listen_count,omitempty"`
		Duration      int64  `json:"duration,omitempty"`
		Lossless      bool   `json:"lossless,omitempty"`
		IsInFavorites bool   `json:"is_in_favorites,omitempty"`
	}

	TrackID struct {
		ID int64 `json:"id,omitempty" form:"id" query:"id"`
	}
)

func (t *TrackAlbum) BindProto(track *proto.AlbumTrack) {
	bindedTrack := &TrackAlbum{
		ID:            track.ID,
		Title:         track.Title,
		Explicit:      track.Explicit,
		Genre:         track.Genre,
		Number:        track.Number,
		File:          track.File,
		ListenCount:   track.ListenCount,
		Duration:      track.Duration,
		Lossless:      track.Lossless,
		IsInFavorites: track.IsInFavorites,
	}

	*t = *bindedTrack
}

func (t *Track) BindProto(track *proto.Track) {
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
			ArtworkColor:   track.Album.ArtworkColor,
		},
		Artist: Artist{
			ID:     track.Artist.ID,
			Name:   track.Artist.Name,
			Avatar: track.Artist.Avatar,
			Video:  track.Artist.Video,
		},
		IsInFavorites: track.IsInFavorites,
	}

	*t = *bindedTrack
}
