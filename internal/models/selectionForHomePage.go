package models

type SelectionForHomePage struct {
	Tracks    []Track    `json:"tracks"`
	Albums    []Album    `json:"albums"`
	Playlists []Playlist `json:"playlists"`
	Artists   []Artist   `json:"artists"`
}
