package models

import "2021_2_LostPointer/internal/microservices/music/proto"

type Album struct {
	ID             int64  `json:"id,omitempty"`
	Title          string `json:"title,omitempty"`
	Year           int64  `json:"year,omitempty"`
	Artist         string `json:"artist,omitempty"`
	Artwork        string `json:"artwork,omitempty"`
	TracksCount    int64  `json:"tracks_count,omitempty"`
	TracksDuration int64  `json:"tracks_duration,omitempty"`
	ArtworkColor   string `json:"artwork_color,omitempty"`
}

func (a *Album) BindProtoAlbum(album *proto.Album) {
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
