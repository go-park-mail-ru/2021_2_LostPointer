package models

import (
	music "2021_2_LostPointer/internal/microservices/music/proto"
	playlists "2021_2_LostPointer/internal/microservices/playlists/proto"
)

type PlaylistID struct {
	ID int64 `json:"id,omitempty"`
}

type PlaylistTrack struct {
	TrackID    int64 `json:"track_id,omitempty" form:"track_id" query:"track_id"`
	PlaylistID int64 `json:"playlist_id,omitempty" form:"playlist_id" query:"playlist_id"`
}

type UserPlaylist struct {
	PlaylistID int64  `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Artwork    string `json:"artwork,omitempty"`
	IsPublic   bool   `json:"is_public,omitempty"`
}

type UserPlaylists struct {
	Playlists []UserPlaylist `json:"playlists,omitempty"`
}

type PlaylistPage struct {
	ID           int64   `json:"id,omitempty"`
	Title        string  `json:"title,omitempty"`
	Artwork      string  `json:"artwork,omitempty"`
	ArtworkColor string  `json:"artwork_color,omitempty"`
	Tracks       []Track `json:"tracks,omitempty"`
	IsPublic     bool    `json:"is_public,omitempty"`
}

type PlaylistArtworkColor struct {
	ArtworkColor string `json:"artwork_color"`
}

func (p *PlaylistPage) BindProto(playlistPage *music.PlaylistPageResponse) {
	bindedTracks := make([]Track, 0)
	for _, track := range playlistPage.Tracks {
		var bindedTrack Track
		bindedTrack.BindProto(track)
		bindedTracks = append(bindedTracks, bindedTrack)
	}

	bindedPlaylistPage := PlaylistPage{
		ID:           playlistPage.PlaylistID,
		Title:        playlistPage.Title,
		Artwork:      playlistPage.Artwork,
		ArtworkColor: playlistPage.ArtworkColor,
		Tracks:       bindedTracks,
		IsPublic:     playlistPage.IsPublic,
	}

	*p = bindedPlaylistPage
}

func (u *UserPlaylists) BindProto(playlists *music.PlaylistsData) {
	bindedPlaylists := make([]UserPlaylist, 0)
	for _, playlist := range playlists.Playlists {
		bindedPlaylist := UserPlaylist{
			PlaylistID: playlist.PlaylistID,
			Title:      playlist.Title,
			Artwork:    playlist.Artwork,
			IsPublic:   playlist.IsPublic,
		}
		bindedPlaylists = append(bindedPlaylists, bindedPlaylist)
	}

	u.Playlists = bindedPlaylists
}

func (p *PlaylistID) BindProto(playlist *playlists.CreatePlaylistResponse) {
	binded := &PlaylistID{
		ID: playlist.PlaylistID,
	}

	*p = *binded
}
