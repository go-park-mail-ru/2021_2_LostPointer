package models

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
)

//easyjson:json
type (
	Albums []Album

	Album struct {
		ID             int64  `json:"id,omitempty"`
		Title          string `json:"title,omitempty"`
		Year           int64  `json:"year,omitempty"`
		Artist         string `json:"artist,omitempty"`
		Artwork        string `json:"artwork,omitempty"`
		TracksCount    int64  `json:"tracks_count,omitempty"`
		TracksDuration int64  `json:"tracks_duration,omitempty"`
		ArtworkColor   string `json:"artwork_color,omitempty"`
	}

	AlbumPage struct {
		ID             int64        `json:"id,omitempty"`
		Title          string       `json:"title,omitempty"`
		Year           int64        `json:"year,omitempty"`
		Artwork        string       `json:"artwork,omitempty"`
		TracksCount    int64        `json:"tracks_count,omitempty"`
		TracksDuration int64        `json:"tracks_duration,omitempty"`
		ArtworkColor   string       `json:"artwork_color,omitempty"`
		Artist         Artist       `json:"artist"`
		Tracks         []TrackAlbum `json:"tracks,omitempty"`
	}
)

func (a *AlbumPage) BindProto(album *proto.AlbumPageResponse) {
	tracks := make([]TrackAlbum, 0)
	for _, t := range album.Tracks {
		var track TrackAlbum
		track.BindProto(t)
		tracks = append(tracks, track)
	}

	bindedAlbum := &AlbumPage{
		ID:             album.AlbumID,
		Title:          album.Title,
		Year:           album.Year,
		Artwork:        album.Artwork,
		TracksCount:    album.TracksCount,
		TracksDuration: album.TracksDuration,
		ArtworkColor:   album.ArtworkColor,
		Artist: Artist{
			ID:   album.Artist.ID,
			Name: album.Artist.Name,
		},
		Tracks: tracks,
	}

	*a = *bindedAlbum
}

func (a *Album) BindProto(album *proto.Album) {
	bindedAlbum := &Album{
		ID:             album.ID,
		Title:          album.Title,
		Year:           album.Year,
		Artist:         album.Artist,
		Artwork:        album.Artwork,
		TracksCount:    album.TracksAmount,
		TracksDuration: album.TracksDuration,
		ArtworkColor:   album.ArtworkColor,
	}

	*a = *bindedAlbum
}
