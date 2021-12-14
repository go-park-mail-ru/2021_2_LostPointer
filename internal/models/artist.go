package models

import "2021_2_LostPointer/internal/microservices/music/proto"

//easyjson:json
type Artists []Artist

type Artist struct {
	ID     int64   `json:"id,omitempty"`
	Name   string  `json:"name"`
	Avatar string  `json:"avatar,omitempty"`
	Video  string  `json:"video,omitempty"`
	Tracks []Track `json:"tracks,omitempty"`
	Albums []Album `json:"albums,omitempty"`
}

func (a *Artist) BindProto(artist *proto.Artist) {
	tracks := make([]Track, 0)
	for _, t := range artist.Tracks {
		var track Track
		track.BindProto(t)
		tracks = append(tracks, track)
	}

	albums := make([]Album, 0)
	for _, alb := range artist.Albums {
		var album Album
		album.BindProto(alb)
		albums = append(albums, album)
	}

	bindedArtists := &Artist{
		ID:     artist.ID,
		Name:   artist.Name,
		Avatar: artist.Avatar,
		Video:  artist.Video,
		Tracks: tracks,
		Albums: albums,
	}

	*a = *bindedArtists
}
