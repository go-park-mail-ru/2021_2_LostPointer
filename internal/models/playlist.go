package models

import "2021_2_LostPointer/internal/microservices/playlists/proto"

type PlaylistID struct {
	ID int64 `json:"id,omitempty"`
}

type PlaylistTrack struct {
	TrackID    int64 `json:"track_id,omitempty" form:"track_id" query:"track_id"`
	PlaylistID int64 `json:"playlist_id,omitempty" form:"playlist_id" query:"playlist_id"`
}

func (p *PlaylistID) BindProto(playlist *proto.CreatePlaylistResponse) {
	binded := &PlaylistID{
		ID: playlist.PlaylistID,
	}

	*p = *binded
}
